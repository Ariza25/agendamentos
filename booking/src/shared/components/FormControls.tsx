import { forwardRef } from 'react'
import type {
  InputHTMLAttributes,
  LabelHTMLAttributes,
  FormHTMLAttributes,
  ReactNode,
  SelectHTMLAttributes,
  TextareaHTMLAttributes,
} from 'react'

type InputProps = InputHTMLAttributes<HTMLInputElement>
type SelectProps = SelectHTMLAttributes<HTMLSelectElement>
type TextareaProps = TextareaHTMLAttributes<HTMLTextAreaElement>
type CheckboxProps = InputHTMLAttributes<HTMLInputElement>
type FormProps = FormHTMLAttributes<HTMLFormElement>

type FieldProps = LabelHTMLAttributes<HTMLLabelElement> & {
  label?: ReactNode
  error?: ReactNode
  hint?: ReactNode
}

const controlClassName =
  'w-full min-w-0 rounded-xl border border-[var(--color-line)] bg-white px-3 text-sm font-bold text-[var(--color-ink)] outline-none transition focus:border-[var(--color-primary)] disabled:cursor-not-allowed disabled:opacity-70'

function mergeClasses(...classes: Array<string | undefined>) {
  return classes.filter(Boolean).join(' ')
}

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input({ className, ...props }, ref) {
  return <input {...props} className={mergeClasses(controlClassName, 'h-11', className)} ref={ref} />
})

export const Select = forwardRef<HTMLSelectElement, SelectProps>(function Select({ className, ...props }, ref) {
  return (
    <select
      {...props}
      className={mergeClasses(controlClassName, 'admin-select-control h-11 pl-3 pr-10', className)}
      ref={ref}
    />
  )
})

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(function Textarea({ className, ...props }, ref) {
  return <textarea {...props} className={mergeClasses(controlClassName, 'min-h-24 py-3 font-medium', className)} ref={ref} />
})

export function Field({ label, error, hint, className, children, ...props }: FieldProps) {
  return (
    <label {...props} className={mergeClasses('block min-w-0', className)}>
      {label ? <span className="text-sm font-bold text-[var(--color-ink)]">{label}</span> : null}
      {children}
      {hint ? <span className="mt-2 block text-xs font-semibold text-[var(--color-ink-soft)]">{hint}</span> : null}
      {error ? <span className="mt-2 block text-xs font-semibold text-red-700">{error}</span> : null}
    </label>
  )
}

export function Form({ className, children, ...props }: FormProps) {
  return (
    <form {...props} className={mergeClasses('mt-5 grid gap-4', className)}>
      {children}
    </form>
  )
}
