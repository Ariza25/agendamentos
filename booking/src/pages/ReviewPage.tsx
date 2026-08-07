import { ArrowLeft, Star, ThumbsUp } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import { fetchPublicReview, submitPublicReview } from '../lib/bookingApi'
import type { PublicReview } from '../lib/bookingApi'
import { Button } from '../shared/components/Button'
import { Field, Form, Textarea } from '../shared/components/FormControls'
import { Spinner } from '../shared/components/Spinner'

export function ReviewPage() {
  const { slug = '', reviewId = '' } = useParams<{ slug: string; reviewId: string }>()
  const navigate = useNavigate()

  const [review, setReview] = useState<PublicReview | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState('')

  const [rating, setRating] = useState(0)
  const [hoveredRating, setHoveredRating] = useState(0)
  const [comment, setComment] = useState('')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isDone, setIsDone] = useState(false)

  useEffect(() => {
    let active = true
    setLoading(true)
    fetchPublicReview(slug, reviewId)
      .then((data) => {
        if (!active) {
          return
        }
        setReview(data)
        if (data.status === 'submitted') {
          setIsDone(true)
          setRating(data.rating)
          setComment(data.comment)
        }
      })
      .catch((loadFail: unknown) => {
        if (active) {
          setLoadError((loadFail as Error).message)
        }
      })
      .finally(() => {
        if (active) {
          setLoading(false)
        }
      })
    return () => {
      active = false
    }
  }, [slug, reviewId])

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (rating < 1) {
      setError('Toque em uma estrela para dar sua nota.')
      return
    }
    setIsSubmitting(true)
    setError('')
    try {
      await submitPublicReview(slug, reviewId, rating, comment.trim())
      setIsDone(true)
    } catch (submitError) {
      setError((submitError as Error).message)
    } finally {
      setIsSubmitting(false)
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-dvh items-center justify-center px-4 py-10">
        <Spinner label="Carregando avaliação..." />
      </main>
    )
  }

  return (
    <main className="flex min-h-dvh items-center justify-center px-4 py-10">
      <section className="w-full max-w-md rounded-3xl border border-[var(--color-line)] bg-white p-8 shadow-sm">
        {loadError ? (
          <div className="text-center">
            <h1 className="font-headline text-xl font-extrabold text-[var(--color-ink)]">Não foi possível carregar</h1>
            <p className="mt-2 text-sm font-semibold text-[var(--color-ink-soft)]">{loadError}</p>
            <Button className="mt-6" onClick={() => navigate('/')} size="lg" variant="secondary">
              Voltar ao início
            </Button>
          </div>
        ) : isDone ? (
          <div className="text-center">
            <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-[var(--color-success)] text-white">
              <ThumbsUp className="h-7 w-7" />
            </span>
            <h1 className="mt-4 font-headline text-2xl font-extrabold text-[var(--color-ink)]">
              {review?.status === 'submitted' ? 'Obrigado pela avaliação!' : 'Avaliação enviada!'}
            </h1>
            <p className="mt-2 text-sm font-semibold text-[var(--color-ink-soft)]">
              {review?.status === 'submitted'
                ? 'Sua opinião já está registrada. Agradecemos o retorno.'
                : 'Sua opinião ajuda a melhorar o atendimento. Muito obrigado!'}
            </p>
            <Button className="mt-6" onClick={() => navigate('/')} size="lg" variant="secondary">
              Voltar ao início
            </Button>
          </div>
        ) : (
          <>
            <div className="flex items-center gap-3">
              <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
                <Star className="h-6 w-6" />
              </span>
              <div>
                <h1 className="font-headline text-2xl font-extrabold text-[var(--color-ink)]">Como foi o atendimento?</h1>
                <p className="text-sm font-semibold text-[var(--color-ink-soft)]">Sua avaliação leva menos de um minuto.</p>
              </div>
            </div>

            {error ? (
              <p className="mt-4 rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p>
            ) : null}

            <Form onSubmit={handleSubmit}>
              <Field label="Sua nota">
                <div className="mt-2 flex items-center gap-1">
                  {[1, 2, 3, 4, 5].map((value) => (
                    <button
                      aria-label={`${value} estrela${value > 1 ? 's' : ''}`}
                      className="rounded-full p-1 transition hover:scale-110"
                      key={value}
                      onClick={() => setRating(value)}
                      onMouseEnter={() => setHoveredRating(value)}
                      onMouseLeave={() => setHoveredRating(0)}
                      type="button"
                    >
                      <Star
                        className={`h-9 w-9 ${
                          value <= (hoveredRating || rating) ? 'fill-amber-400 text-amber-400' : 'text-stone-300'
                        }`}
                      />
                    </button>
                  ))}
                </div>
              </Field>
              <Field hint="Opcional" label="Comentário">
                <Textarea
                  className="mt-2"
                  maxLength={500}
                  onChange={(event) => setComment(event.target.value)}
                  placeholder="Conte como foi sua experiência..."
                  value={comment}
                />
              </Field>
              <Button icon={<ThumbsUp className="h-4 w-4" />} loading={isSubmitting} size="lg" type="submit">
                Enviar avaliação
              </Button>
              <Button
                className="mt-2 w-full"
                icon={<ArrowLeft className="h-4 w-4" />}
                onClick={() => navigate('/')}
                variant="ghost"
              >
                Voltar ao início
              </Button>
            </Form>
          </>
        )}
      </section>
    </main>
  )
}
