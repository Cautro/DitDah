import { useEffect, useMemo, useState } from 'react'

type Language = 'ru' | 'en'
type TaskType = 'text' | 'tap' | 'listen' | 'quiz'

type User = {
  id: number
  username: string
  isAdmin: boolean
}

type TextPayload = {
  text: string
  morse?: string
}

type QuizPayload = {
  question: string
  options: string[]
  correctIndex: number
}

type ListenPayload = {
  morseText: string
  correctAnswer: string
}

type TapPayload = {
  words: string
  morse: string
}

type TaskPayload = TextPayload | QuizPayload | ListenPayload | TapPayload

type Task = {
  id: number
  lessonId: number
  order: number
  type: TaskType
  payload: TaskPayload
}

type LessonData = {
  id: number
  order: number
  title: string
  description: string
  xpReward: number
  language: Language
  tasks: Task[]
}

type FrontendTask = Task & {
  _frontendId: string
}

type LessonDraft = Omit<LessonData, 'tasks'> & {
  tasks: FrontendTask[]
}

const taskTypeLabels: Record<TaskType, string> = {
  text: 'Text',
  tap: 'Tap',
  listen: 'Listen',
  quiz: 'Quiz',
}

const makeFrontendId = () =>
  typeof crypto !== 'undefined' && 'randomUUID' in crypto
    ? crypto.randomUUID()
    : `${Date.now()}-${Math.random().toString(16).slice(2)}`

const emptyQuizOptions = () => ['', '', '', '']

const createTask = (type: TaskType, order: number): FrontendTask => {
  switch (type) {
    case 'text':
      return {
        _frontendId: makeFrontendId(),
        id: 0,
        lessonId: 0,
        order,
        type,
        payload: {
          text: 'Сигнал для распознавания',
          morse: '.-.-.',
        },
      }
    case 'tap':
      return {
        _frontendId: makeFrontendId(),
        id: 0,
        lessonId: 0,
        order,
        type,
        payload: {
          words: 'SOS',
          morse: '... --- ...',
        },
      }
    case 'listen':
      return {
        _frontendId: makeFrontendId(),
        id: 0,
        lessonId: 0,
        order,
        type,
        payload: {
          morseText: '.... . .-.. .-.. ---',
          correctAnswer: 'HELLO',
        },
      }
    case 'quiz':
      return {
        _frontendId: makeFrontendId(),
        id: 0,
        lessonId: 0,
        order,
        type,
        payload: {
          question: 'Что означает «... --- ...»?',
          options: ['SOS', 'OK', 'HI', 'YES'],
          correctIndex: 0,
        },
      }
    default:
      return {
        _frontendId: makeFrontendId(),
        id: 0,
        lessonId: 0,
        order,
        type: 'text',
        payload: {
          text: '',
          morse: '',
        },
      }
  }
}

const initialDraft = (): LessonDraft => ({
  id: 0,
  order: 1,
  title: 'Основы азбуки Морзе',
  description: 'Освой базовое распознавание символов, пауз и простых сигналов.',
  xpReward: 120,
  language: 'ru',
  tasks: [
    createTask('text', 1),
    createTask('tap', 2),
    createTask('listen', 3),
    createTask('quiz', 4),
    createTask('text', 5),
    createTask('tap', 6),
    createTask('quiz', 7),
  ],
})

const toBackendLesson = (draft: LessonDraft) => {
  return {
    order: Number(draft.order) || 0,
    title: draft.title.trim(),
    description: draft.description.trim(),
    xpReward: Number(draft.xpReward) || 0,
    language: draft.language,
    tasks: draft.tasks.map((task, index) => {
      const { _frontendId, id: _taskId, lessonId: _lessonId, ...taskData } = task

      return {
        ...taskData,
        order: index,
        payload: normalizePayloadForBackend(task.type, task.payload),
      }
    }),
  }
}

const normalizePayloadForBackend = (type: TaskType, payload: TaskPayload): TaskPayload => {
  switch (type) {
    case 'text': {
      const text = payload as TextPayload
      return {
        text: text.text,
        morse: text.morse ?? '',
      }
    }
    case 'tap': {
      const tap = payload as TapPayload
      return {
        words: tap.words,
        morse: tap.morse,
      }
    }
    case 'listen': {
      const listen = payload as ListenPayload
      return {
        morseText: listen.morseText,
        correctAnswer: listen.correctAnswer,
      }
    }
    case 'quiz': {
      const quiz = payload as QuizPayload
      const options = [...(quiz.options ?? emptyQuizOptions())].map((option) => option ?? '')
      const validIndex =
        Number.isInteger(quiz.correctIndex) && quiz.correctIndex >= 0 && quiz.correctIndex < options.length
          ? quiz.correctIndex
          : 0

      return {
        question: quiz.question,
        options,
        correctIndex: validIndex,
      }
    }
    default:
      return payload
  }
}

async function apiFetch<T>(url: string, init: RequestInit = {}, retry = true): Promise<T> {
  const response = await fetch(url, {
    credentials: 'include',
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers ?? {}),
    },
  })

  if (response.status === 401 && retry) {
    const refreshRes = await fetch('/refresh', {
      method: 'POST',
      credentials: 'include',
    })

    if (refreshRes.ok) {
      return apiFetch<T>(url, init, false)
    }

    throw new Error('Unauthorized')
  }

  if (!response.ok) {
    let message = `HTTP ${response.status}`

    try {
      const data = (await response.json()) as { error?: string; message?: string }
      message = data.error ?? data.message ?? message
    } catch {
      // ignore
    }

    throw new Error(message)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

function AdminLessonBuilder() {
  const [lesson, setLesson] = useState<LessonDraft>(() => initialDraft())
  const [status, setStatus] = useState<string>('')
  const [token, setToken] = useState('')

  const taskCount = lesson.tasks.length
  const isValid =
    lesson.title.trim().length > 0 &&
    lesson.description.trim().length > 0 &&
    taskCount >= 7 &&
    taskCount <= 20

  const preview = useMemo(() => JSON.stringify(toBackendLesson(lesson), null, 2), [lesson])

  const updateLessonField = <K extends keyof LessonDraft>(key: K, value: LessonDraft[K]) => {
    setLesson((prev) => ({ ...prev, [key]: value }))
  }

  const addTask = (type: TaskType) => {
    setLesson((prev) => ({
      ...prev,
      tasks: [...prev.tasks, createTask(type, prev.tasks.length + 1)],
    }))
  }

  const removeTask = (frontendId: string) => {
    setLesson((prev) => ({
      ...prev,
      tasks: prev.tasks.filter((task) => task._frontendId !== frontendId),
    }))
  }

  const updateTask = (frontendId: string, updates: Partial<FrontendTask>) => {
    setLesson((prev) => ({
      ...prev,
      tasks: prev.tasks.map((task) =>
        task._frontendId === frontendId ? { ...task, ...updates } : task,
      ),
    }))
  }

  const moveTask = (index: number, direction: -1 | 1) => {
    const targetIndex = index + direction
    if (targetIndex < 0 || targetIndex >= lesson.tasks.length) {
      return
    }

    setLesson((prev) => {
      const nextTasks = [...prev.tasks]
      const [current] = nextTasks.splice(index, 1)
      nextTasks.splice(targetIndex, 0, current)

      return {
        ...prev,
        tasks: nextTasks.map((task, orderIndex) => ({
          ...task,
          order: orderIndex + 1,
        })),
      }
    })
  }

  const updateTaskPayload = (frontendId: string, payload: Partial<Record<string, unknown>>) => {
    setLesson((prev) => ({
      ...prev,
      tasks: prev.tasks.map((task) => {
        if (task._frontendId !== frontendId) {
          return task
        }

        return {
          ...task,
          payload: {
            ...task.payload,
            ...payload,
          },
        }
      }),
    }))
  }

  const updateQuizOption = (frontendId: string, index: number, value: string) => {
    setLesson((prev) => ({
      ...prev,
      tasks: prev.tasks.map((task) => {
        if (task._frontendId !== frontendId || task.type !== 'quiz') {
          return task
        }

        const quizPayload = task.payload as QuizPayload
        const nextOptions = [...(quizPayload.options ?? emptyQuizOptions())]
        nextOptions[index] = value

        return {
          ...task,
          payload: {
            ...quizPayload,
            options: nextOptions,
          },
        }
      }),
    }))
  }

  const addQuizOption = (frontendId: string) => {
    setLesson((prev) => ({
      ...prev,
      tasks: prev.tasks.map((task) => {
        if (task._frontendId !== frontendId || task.type !== 'quiz') {
          return task
        }

        const quizPayload = task.payload as QuizPayload
        const currentOptions = [...(quizPayload.options ?? emptyQuizOptions())]
        currentOptions.push('')

        return {
          ...task,
          payload: {
            ...quizPayload,
            options: currentOptions,
          },
        }
      }),
    }))
  }

  const removeQuizOption = (frontendId: string, optionIndex: number) => {
    setLesson((prev) => ({
      ...prev,
      tasks: prev.tasks.map((task) => {
        if (task._frontendId !== frontendId || task.type !== 'quiz') {
          return task
        }

        const quizPayload = task.payload as QuizPayload
        const currentOptions = [...(quizPayload.options ?? emptyQuizOptions())]
        if (currentOptions.length <= 2) return task

        currentOptions.splice(optionIndex, 1)

        return {
          ...task,
          payload: {
            ...quizPayload,
            options: currentOptions,
            correctIndex:
              quizPayload.correctIndex >= currentOptions.length
                ? currentOptions.length - 1
                : quizPayload.correctIndex,
          },
        }
      }),
    }))
  }

  const handleSubmit = async () => {
    if (!isValid) {
      setStatus('Заполните title, description и убедитесь, что у урока от 7 до 20 задач.')
      return
    }

    const payload = toBackendLesson(lesson)

    try {
      await apiFetch('/api/lesson/add', {
        method: 'POST',
        body: JSON.stringify(payload),
        headers: {
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      })

      setStatus('Урок успешно отправлен в backend')
    } catch (error) {
      setStatus(`Ошибка отправки: ${error instanceof Error ? error.message : 'unknown error'}`)
    }
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="mx-auto max-w-7xl px-4 py-6 lg:px-8">
        <div className="mb-6 flex flex-col gap-4 rounded-2xl border border-slate-700 bg-slate-900 p-5 shadow-sm lg:flex-row lg:items-center lg:justify-between">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-indigo-400">Lesson Builder</p>
            <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-50">Конструктор урока</h1>
          </div>

          <button
            type="button"
            onClick={handleSubmit}
            disabled={!isValid}
            className="rounded-xl bg-indigo-600 px-5 py-3 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-500 disabled:cursor-not-allowed disabled:bg-slate-700"
          >
            Сохранить / Отправить
          </button>
        </div>

        {!isValid && (
          <div className="mb-6 rounded-xl border border-amber-700 bg-amber-950/70 px-4 py-3 text-sm text-amber-200">
            Требуется заполнить title, description и количество задач от 7 до 20.
          </div>
        )}

        <div className="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
          <aside className="rounded-2xl border border-slate-700 bg-slate-900 p-5 shadow-sm">
            <div className="mb-5 space-y-1">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">Lesson</p>
              <h2 className="text-xl font-bold text-slate-50">{lesson.title || 'Новый урок'}</h2>
            </div>

            <div className="mb-5 rounded-xl bg-slate-800 p-4">
              <div className="flex items-center justify-between text-sm text-slate-400">
                <span>Language</span>
                <span className="font-semibold uppercase text-slate-100">{lesson.language}</span>
              </div>
              <div className="mt-3 flex items-center justify-between text-sm text-slate-400">
                <span>XP</span>
                <span className="font-semibold text-slate-100">{lesson.xpReward}</span>
              </div>
              <div className="mt-3 flex items-center justify-between text-sm text-slate-400">
                <span>Tasks</span>
                <span className="font-semibold text-slate-100">{taskCount}</span>
              </div>
            </div>

            <div className="space-y-2">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">Добавить задание</p>
              <div className="grid grid-cols-2 gap-2">
                {(['text', 'tap', 'listen', 'quiz'] as TaskType[]).map((type) => (
                  <button key={type} type="button" onClick={() => addTask(type)} className="rounded-xl border border-slate-700 bg-slate-800 px-3 py-2 text-sm font-medium text-slate-200 shadow-sm transition hover:border-indigo-500 hover:text-indigo-300">
                    {taskTypeLabels[type]}
                  </button>
                ))}
              </div>
            </div>
          </aside>

          <main className="space-y-6">
            <section className="rounded-2xl border border-slate-700 bg-slate-900 p-5 shadow-sm">
              <h3 className="mb-4 text-lg font-semibold text-slate-50">Метаданные урока</h3>

              <div className="grid gap-4 md:grid-cols-2">
                <label className="space-y-2 text-sm font-medium text-slate-300">
                  <span>Title</span>
                  <input
                    value={lesson.title}
                    onChange={(event) => updateLessonField('title', event.target.value)}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition placeholder:text-slate-500 focus:border-indigo-500 focus:bg-slate-900"
                    placeholder="Например: Основы азбуки Морзе"
                  />
                </label>

                <label className="space-y-2 text-sm font-medium text-slate-300">
                  <span>Language</span>
                  <select
                    value={lesson.language}
                    onChange={(event) => updateLessonField('language', event.target.value as Language)}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500 focus:bg-slate-900"
                  >
                    <option value="ru">ru</option>
                    <option value="en">en</option>
                  </select>
                </label>

                <label className="space-y-2 text-sm font-medium text-slate-300 md:col-span-2">
                  <span>Description</span>
                  <textarea
                    value={lesson.description}
                    onChange={(event) => updateLessonField('description', event.target.value)}
                    rows={3}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition placeholder:text-slate-500 focus:border-indigo-500 focus:bg-slate-900"
                    placeholder="Кратко опишите, чему научится пользователь"
                  />
                </label>

                <label className="space-y-2 text-sm font-medium text-slate-300">
                  <span>Order</span>
                  <input
                    type="number"
                    min={1}
                    value={lesson.order}
                    onChange={(event) => updateLessonField('order', Number(event.target.value) || 1)}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500 focus:bg-slate-900"
                  />
                </label>

                <label className="space-y-2 text-sm font-medium text-slate-300">
                  <span>XP Reward</span>
                  <input
                    type="number"
                    min={0}
                    value={lesson.xpReward}
                    onChange={(event) => updateLessonField('xpReward', Number(event.target.value) || 0)}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500 focus:bg-slate-900"
                  />
                </label>
              </div>
            </section>

            <section className="rounded-2xl border border-slate-700 bg-slate-900 p-5 shadow-sm">
              <div className="mb-4 flex items-center justify-between gap-3">
                <h3 className="text-lg font-semibold text-slate-50">Задания</h3>
                <span className={`rounded-full px-2.5 py-1 text-xs font-semibold ${taskCount >= 7 && taskCount <= 20 ? 'bg-emerald-900/80 text-emerald-200' : 'bg-amber-900/80 text-amber-200'}`}>
                  {taskCount} / 7-20
                </span>
              </div>

              <div className="space-y-4">
                {lesson.tasks.map((task, index) => (
                  <div key={task._frontendId} className="rounded-2xl border border-slate-700 bg-slate-800 p-4">
                    <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                      <div className="flex items-center gap-3">
                        <span className="inline-flex h-8 w-8 items-center justify-center rounded-xl bg-indigo-900/80 font-semibold text-indigo-200">
                          {index + 1}
                        </span>
                        <span className="text-sm font-semibold uppercase tracking-[0.12em] text-slate-300">
                          {taskTypeLabels[task.type]}
                        </span>
                      </div>

                      <div className="flex items-center gap-2">
                        <button
                          type="button"
                          onClick={() => moveTask(index, -1)}
                          disabled={index === 0}
                          className="rounded-lg border border-slate-600 bg-slate-900 px-2 py-1 text-sm text-slate-300 disabled:cursor-not-allowed disabled:opacity-40"
                        >
                          ↑
                        </button>
                        <button
                          type="button"
                          onClick={() => moveTask(index, 1)}
                          disabled={index === lesson.tasks.length - 1}
                          className="rounded-lg border border-slate-600 bg-slate-900 px-2 py-1 text-sm text-slate-300 disabled:cursor-not-allowed disabled:opacity-40"
                        >
                          ↓
                        </button>
                        <button
                          type="button"
                          onClick={() => removeTask(task._frontendId)}
                          className="rounded-lg border border-rose-700 bg-rose-950/60 px-2 py-1 text-sm font-medium text-rose-200"
                        >
                          Удалить
                        </button>
                      </div>
                    </div>

                    {task.type === 'text' && (
                      <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-2 text-sm font-medium text-slate-300">
                          <span>Text</span>
                          <textarea
                            rows={3}
                            value={(task.payload as TextPayload).text ?? ''}
                            onChange={(event) =>
                              updateTaskPayload(task._frontendId, { text: event.target.value })
                            }
                            className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                          />
                        </label>

                        <label className="space-y-2 text-sm font-medium text-slate-300">
                          <span>Morse</span>
                          <input
                            value={(task.payload as TextPayload).morse ?? ''}
                            onChange={(event) =>
                              updateTaskPayload(task._frontendId, { morse: event.target.value })
                            }
                            className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                          />
                        </label>
                      </div>
                    )}

                    {task.type === 'tap' && (
                      <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-2 text-sm font-medium text-slate-300">
                          <span>Words</span>
                          <input
                            value={(task.payload as TapPayload).words ?? ''}
                            onChange={(event) =>
                              updateTaskPayload(task._frontendId, { words: event.target.value })
                            }
                            className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                          />
                        </label>

                        <label className="space-y-2 text-sm font-medium text-slate-300">
                          <span>Morse</span>
                          <input
                            value={(task.payload as TapPayload).morse ?? ''}
                            onChange={(event) =>
                              updateTaskPayload(task._frontendId, { morse: event.target.value })
                            }
                            className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                          />
                        </label>
                      </div>
                    )}

                    {task.type === 'listen' && (
                      <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-2 text-sm font-medium text-slate-300">
                          <span>morseText</span>
                          <input
                            value={(task.payload as ListenPayload).morseText ?? ''}
                            onChange={(event) =>
                              updateTaskPayload(task._frontendId, { morseText: event.target.value })
                            }
                            className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                          />
                        </label>

                        <label className="space-y-2 text-sm font-medium text-slate-300">
                          <span>correctAnswer</span>
                          <input
                            value={(task.payload as ListenPayload).correctAnswer ?? ''}
                            onChange={(event) =>
                              updateTaskPayload(task._frontendId, { correctAnswer: event.target.value })
                            }
                            className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                          />
                        </label>
                      </div>
                    )}

                    {task.type === 'quiz' && (
                      <div className="space-y-4">
                        <label className="block space-y-2 text-sm font-medium text-slate-300">
                          <span>Question</span>
                          <input
                            value={(task.payload as QuizPayload).question ?? ''}
                            onChange={(event) =>
                              updateTaskPayload(task._frontendId, { question: event.target.value })
                            }
                            className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                          />
                        </label>

                        <div className="space-y-3">
                          {((task.payload as QuizPayload).options ?? emptyQuizOptions()).map((option, optionIndex) => (
                            <div key={`${task._frontendId}-option-${optionIndex}`} className="flex items-center gap-3">
                              <input
                                type="radio"
                                name={`correct-${task._frontendId}`}
                                checked={(task.payload as QuizPayload).correctIndex === optionIndex}
                                onChange={() =>
                                  updateTask(task._frontendId, {
                                    payload: {
                                      ...task.payload,
                                      correctIndex: optionIndex,
                                    },
                                  })
                                }
                                className="h-4 w-4 accent-indigo-500"
                              />
                              <input
                                value={option}
                                onChange={(event) =>
                                  updateQuizOption(task._frontendId, optionIndex, event.target.value)
                                }
                                className="w-full rounded-xl border border-slate-700 bg-slate-900 px-3 py-2.5 text-slate-50 outline-none transition focus:border-indigo-500"
                                placeholder={`Option ${optionIndex + 1}`}
                              />
                              <button
                                type="button"
                                onClick={() => removeQuizOption(task._frontendId, optionIndex)}
                                disabled={((task.payload as QuizPayload).options?.length ?? 0) <= 2}
                                className="rounded-lg border border-slate-600 bg-slate-900 px-2 py-2 text-xs text-slate-300 disabled:cursor-not-allowed disabled:opacity-40"
                              >
                                Удалить
                              </button>
                            </div>
                          ))}
                        </div>

                        <button
                          type="button"
                          onClick={() => addQuizOption(task._frontendId)}
                          className="rounded-xl border border-dashed border-indigo-500 bg-indigo-950/50 px-3 py-2 text-sm font-medium text-indigo-200"
                        >
                          + Добавить вариант
                        </button>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </section>

            <section className="rounded-2xl border border-slate-700 bg-slate-900 p-5 shadow-sm">
              <div className="mb-3 flex items-center justify-between gap-3">
                <h3 className="text-lg font-semibold text-slate-50">JSON payload</h3>
                <label className="text-sm text-slate-400">JWT token</label>
              </div>

              <input
                value={token}
                onChange={(event) => setToken(event.target.value)}
                placeholder="Bearer token (необязательно)"
                className="mb-4 w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition placeholder:text-slate-500 focus:border-indigo-500 focus:bg-slate-900"
              />

              <pre className="overflow-x-auto rounded-2xl bg-slate-950 p-4 text-sm leading-6 text-slate-200">{preview}</pre>

              {status && (
                <div className="mt-4 rounded-xl border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-200">
                  {status}
                </div>
              )}
            </section>
          </main>
        </div>
      </div>
    </div>
  )
}

function LoginScreen({
  onLoggedIn,
}: {
  onLoggedIn: (user: User) => void
}) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    setLoading(true)
    setError('')

    try {
      await apiFetch('/login', {
        method: 'POST',
        body: JSON.stringify({ username, password }),
      })

      const user = await apiFetch<User>('/api/me')
      onLoggedIn(user)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-950 px-4 text-slate-100">
      <div className="w-full max-w-md rounded-2xl border border-slate-700 bg-slate-900 p-8 shadow-lg">
        <div className="mb-6">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-indigo-400">DitDah</p>
          <h1 className="mt-2 text-3xl font-bold text-slate-50">Авторизация</h1>
        </div>

        <form onSubmit={submit} className="space-y-4">
          <label className="block space-y-2 text-sm font-medium text-slate-300">
            <span>Username</span>
            <input
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              className="w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition placeholder:text-slate-500 focus:border-indigo-500 focus:bg-slate-900"
              placeholder="admin"
            />
          </label>

          <label className="block space-y-2 text-sm font-medium text-slate-300">
            <span>Password</span>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="w-full rounded-xl border border-slate-700 bg-slate-800 px-3 py-2.5 text-slate-50 outline-none transition placeholder:text-slate-500 focus:border-indigo-500 focus:bg-slate-900"
              placeholder="••••••••"
            />
          </label>

          {error && (
            <div className="rounded-xl border border-rose-800 bg-rose-950/60 px-3 py-2 text-sm text-rose-200">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-xl bg-indigo-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-indigo-500 disabled:cursor-not-allowed disabled:bg-slate-700"
          >
            {loading ? 'Вход...' : 'Войти'}
          </button>
        </form>
      </div>
    </div>
  )
}

function AccessDenied({
  user,
  onLogout,
}: {
  user: User
  onLogout: () => void
}) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-950 px-4">
      <div className="max-w-lg rounded-2xl border border-amber-700 bg-slate-900 p-8 text-center shadow-lg">
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-amber-400">Access denied</p>
        <h1 className="mt-3 text-3xl font-bold text-slate-50">Доступ только для admin</h1>
        <p className="mt-3 text-slate-300">
          Пользователь <span className="font-semibold text-slate-50">{user.username}</span> не имеет прав администратора.
        </p>

        <button
          type="button"
          onClick={onLogout}
          className="mt-6 rounded-xl bg-slate-100 px-4 py-3 text-sm font-semibold text-slate-900 transition hover:bg-slate-200"
        >
          Выйти
        </button>
      </div>
    </div>
  )
}

function App() {
  const [user, setUser] = useState<User | null>(null)
  const [checking, setChecking] = useState(true)

  const loadSession = async () => {
    try {
      const me = await apiFetch<User>('/api/me')
      setUser(me)
    } catch {
      setUser(null)
    } finally {
      setChecking(false)
    }
  }

  useEffect(() => {
    void loadSession()
  }, [])

  const handleLogout = async () => {
    try {
      await apiFetch('/api/logout', { method: 'POST' })
    } catch {
      // ignore
    } finally {
      setUser(null)
      setChecking(false)
    }
  }

  if (checking) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-950 text-slate-300">
        Проверка авторизации...
      </div>
    )
  }

  if (!user) {
    return <LoginScreen onLoggedIn={(nextUser) => setUser(nextUser)} />
  }

  if (!user.isAdmin) {
    return <AccessDenied user={user} onLogout={handleLogout} />
  }

  return <AdminLessonBuilder />
}

export default App
