import { test, expect } from '@playwright/test';

test.describe('Авторизация', () => {
  test('Успешный вход', async ({ page }) => {
    await page.goto('/');
    await page.fill('#login', 'testuser');
    await page.fill('#password', 'test123');
    await page.click('button:has-text("Войти")');
    await page.waitForURL('/', { timeout: 10000 });
    await expect(page).toHaveURL('/');
  });

  test('Ошибка при неверном пароле', async ({ page }) => {
    await page.goto('/');
    await page.fill('#login', 'testuser');
    await page.fill('#password', 'wrongpassword');
    await page.click('button:has-text("Войти")');
    
    // Проверяем, что появился текст ошибки
    await expect(page.locator('text=invalid login or password')).toBeVisible({ timeout: 5000 });
  });
});
