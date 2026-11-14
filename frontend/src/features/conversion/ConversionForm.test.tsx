import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import { ConversionForm } from './ConversionForm'

vi.mock('@/shared/api/temperatureApi', () => ({
  postConvert: vi.fn(),
}))

import { postConvert } from '@/shared/api/temperatureApi'

it('shows result after successful conversion', async () => {
  vi.mocked(postConvert).mockResolvedValue({ result: 32 })
  render(<ConversionForm />)
  await userEvent.type(screen.getByLabelText('Valor'), '0')
  await userEvent.click(screen.getByRole('button', { name: /convertir/i }))
  expect(await screen.findByText(/32/)).toBeInTheDocument()
})

it('shows error when API fails', async () => {
  vi.mocked(postConvert).mockRejectedValue(new Error('invalid unit: "bad"'))
  render(<ConversionForm />)
  await userEvent.type(screen.getByLabelText('Valor'), '0')
  await userEvent.click(screen.getByRole('button', { name: /convertir/i }))
  expect(await screen.findByRole('alert')).toBeInTheDocument()
})
