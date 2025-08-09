import { test, expect } from '@playwright/test'

test.describe('Targets Management', () => {
  test.beforeEach(async ({ page }) => {
    // Mock initial targets API call
    await page.route('/api/targets', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            id: 1,
            name: 'Production DB',
            host: 'prod.example.com',
            port: 3306,
            db_name: 'production',
            user: 'backup_user',
            comment: 'Production database backup',
            schedule_time: '02:00',
            retention_days: 30,
            auto_compress: true,
            created_at: '2023-01-01T00:00:00Z',
            updated_at: '2023-01-01T00:00:00Z'
          }
        ])
      })
    })
  })

  test('should display targets list', async ({ page }) => {
    await page.goto('/targets')
    
    await expect(page.locator('h1')).toContainText('Backup Targets')
    await expect(page.locator('.card-title')).toContainText('Production DB')
    await expect(page.locator('text=prod.example.com:3306')).toBeVisible()
    await expect(page.locator('text=Daily at 02:00 UTC')).toBeVisible()
  })

  test('should navigate to add target form', async ({ page }) => {
    await page.goto('/targets')
    await page.click('text=Add Target')
    
    await expect(page).toHaveURL('/targets/new')
    await expect(page.locator('h1')).toContainText('Add Target')
  })

  test('should create a new target', async ({ page }) => {
    // Mock the create API call
    await page.route('/api/targets', async route => {
      if (route.request().method() === 'POST') {
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 2,
            name: 'Test DB',
            host: 'localhost',
            port: 3306,
            db_name: 'testdb',
            user: 'testuser',
            comment: 'Test database',
            schedule_time: '',
            retention_days: 30,
            auto_compress: true,
            created_at: '2023-01-01T01:00:00Z',
            updated_at: '2023-01-01T01:00:00Z'
          })
        })
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([])
        })
      }
    })

    await page.goto('/targets/new')
    
    // Fill out the form
    await page.fill('input[placeholder="My Database"]', 'Test DB')
    await page.fill('input[placeholder="localhost"]', 'localhost')
    await page.fill('input[placeholder="3306"]', '3306')
    await page.fill('input[placeholder="mydb"]', 'testdb')
    await page.fill('input[placeholder="root"]', 'testuser')
    await page.fill('input[type="password"]', 'password123')
    await page.fill('textarea[placeholder="Optional description"]', 'Test database')
    
    // Submit the form
    await page.click('text=Create Target')
    
    // Should redirect to targets list
    await expect(page).toHaveURL('/targets')
    
    // Should show success toast
    await expect(page.locator('.alert-success')).toBeVisible()
  })

  test('should validate required fields', async ({ page }) => {
    await page.goto('/targets/new')
    
    // Try to submit without filling required fields
    await page.click('text=Create Target')
    
    // Should show validation errors
    await expect(page.locator('input[placeholder="My Database"]:invalid')).toBeVisible()
    await expect(page.locator('input[placeholder="localhost"]:invalid')).toBeVisible()
  })

  test('should edit an existing target', async ({ page }) => {
    // Mock the get target API call
    await page.route('/api/targets/1', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 1,
          name: 'Production DB',
          host: 'prod.example.com',
          port: 3306,
          db_name: 'production',
          user: 'backup_user',
          comment: 'Production database backup',
          schedule_time: '02:00',
          retention_days: 30,
          auto_compress: true,
          created_at: '2023-01-01T00:00:00Z',
          updated_at: '2023-01-01T00:00:00Z'
        })
      })
    })

    // Mock the update API call
    await page.route('/api/targets/1', async route => {
      if (route.request().method() === 'PUT') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 1,
            name: 'Updated Production DB',
            host: 'prod.example.com',
            port: 3306,
            db_name: 'production',
            user: 'backup_user',
            comment: 'Updated production database backup',
            schedule_time: '03:00',
            retention_days: 45,
            auto_compress: true,
            created_at: '2023-01-01T00:00:00Z',
            updated_at: '2023-01-01T02:00:00Z'
          })
        })
      }
    })

    await page.goto('/targets/1/edit')
    
    // Form should be pre-filled
    await expect(page.locator('input[value="Production DB"]')).toBeVisible()
    await expect(page.locator('input[value="prod.example.com"]')).toBeVisible()
    
    // Update some fields
    await page.fill('input[value="Production DB"]', 'Updated Production DB')
    await page.fill('textarea', 'Updated production database backup')
    await page.fill('input[value="02:00"]', '03:00')
    await page.fill('input[value="30"]', '45')
    
    // Submit the form
    await page.click('text=Update Target')
    
    // Should redirect to targets list
    await expect(page).toHaveURL('/targets')
  })

  test('should delete a target', async ({ page }) => {
    // Mock the delete API call
    await page.route('/api/targets/1', async route => {
      if (route.request().method() === 'DELETE') {
        await route.fulfill({ status: 200 })
      }
    })

    await page.goto('/targets')
    
    // Click the dropdown menu
    await page.click('[role="button"]:has(svg)')
    
    // Click delete
    await page.click('text=Delete')
    
    // Should show confirmation modal
    await expect(page.locator('.modal')).toBeVisible()
    await expect(page.locator('text=Confirm Deletion')).toBeVisible()
    
    // Confirm deletion
    await page.click('button:has-text("Delete")')
    
    // Should show success toast
    await expect(page.locator('.alert-success')).toBeVisible()
  })

  test('should create manual backup', async ({ page }) => {
    // Mock the backup API call
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

    await page.goto('/targets')
    
    // Click the dropdown menu
    await page.click('[role="button"]:has(svg)')
    
    // Click create backup
    await page.click('text=Create Backup')
    
    // Should show success toast
    await expect(page.locator('.alert-success')).toBeVisible()
    await expect(page.locator('.alert-success')).toContainText('Backup Started')
  })
})

test.describe('Targets with empty state', () => {
  test.beforeEach(async ({ page }) => {
    await page.route('/api/targets', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([])
      })
    })
  })

  test('should show empty state', async ({ page }) => {
    await page.goto('/targets')
    
    await expect(page.locator('text=No targets configured')).toBeVisible()
    await expect(page.locator('text=Get started by adding your first backup target')).toBeVisible()
    
    // Should have Add Target button in empty state
    const addTargetButtons = page.locator('text=Add Target')
    await expect(addTargetButtons).toHaveCount(2) // One in header, one in empty state
  })
})