import { toast } from '../toast'

test('toast is exported', () => {
  expect(toast).toBeDefined()
  expect(typeof toast.success).toBe('function')
  expect(typeof toast.error).toBe('function')
})
