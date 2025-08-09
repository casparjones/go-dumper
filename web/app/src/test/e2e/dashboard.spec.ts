import { test, expect } from '@playwright/test'

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Mock API responses
    await page.route('/api/targets', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            id: 1,
            name: 'Test Database',
            host: 'localhost',
            port: 3306,
            db_name: 'testdb',
            user: 'testuser',
            comment: 'Test database for E2E',
            schedule_time: '02:00',
            retention_days: 30,
            auto_compress: true,
            created_at: '2023-01-01T00:00:00Z',
            updated_at: '2023-01-01T00:00:00Z'
          }
        ])
      })
    })

    await page.goto('/')
  })

  test('should display dashboard title and stats', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Dashboard')
    await expect(page.locator('.stat-title')).toContainText('Total Targets')
    await expect(page.locator('.stat-value')).toContainText('1')
  })

  test('should display recent targets', async ({ page }) => {
    await expect(page.locator('.card-title')).toContainText('Recent Targets')
    await expect(page.locator('.font-semibold')).toContainText('Test Database')
    await expect(page.locator('.text-sm')).toContainText('localhost:3306')
  })

  test('should navigate to add target', async ({ page }) => {
    await page.click('text=Add Target')
    await expect(page).toHaveURL('/targets/new')
  })

  test('should show refresh functionality', async ({ page }) => {
    const refreshButton = page.locator('text=Refresh')
    await expect(refreshButton).toBeVisible()
    await refreshButton.click()
    
    // Should show loading state briefly
    await expect(page.locator('.loading-spinner')).toBeVisible()
  })

  test('should display quick actions', async ({ page }) => {
    const quickActions = page.locator('.card-title:has-text("Quick Actions")')
    await expect(quickActions).toBeVisible()
    
    const addTargetLink = page.locator('a:has-text("Add Target")')
    await expect(addTargetLink).toBeVisible()
    await expect(addTargetLink).toHaveAttribute('href', '/targets/new')
  })

  test('should handle backup action', async ({ page }) => {
    // Mock backup API call
    await page.route('/api/targets/1/backup', async route => {
      await route.fulfill({
        status: 202,
        contentType: 'application/json',
        body: JSON.stringify({
          message: 'Backup started',
          backup_id: 123,
          status: 'running'
        })
      })
    })

    await page.click('text=Backup')
    
    // Should show success toast
    await expect(page.locator('.alert-success')).toBeVisible()
    await expect(page.locator('.alert-success')).toContainText('Backup Started')
  })

  test('should navigate to target backups', async ({ page }) => {
    await page.click('text=View')
    await expect(page).toHaveURL('/targets/1/backups')
  })
})

test.describe('Dashboard with no targets', () => {
  test.beforeEach(async ({ page }) => {
    await page.route('/api/targets', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([])
      })
    })

    await page.goto('/')
  })

  test('should show empty state', async ({ page }) => {
    await expect(page.locator('.stat-value')).toContainText('0')
    await expect(page.locator('text=No targets configured')).toBeVisible()
  })
})

test.describe('Dashboard error handling', () => {
  test('should handle API errors gracefully', async ({ page }) => {
    await page.route('/api/targets', async route => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Server error' })
      })
    })

    await page.goto('/')
    
    // Should show error toast
    await expect(page.locator('.alert-error')).toBeVisible()
  })
})