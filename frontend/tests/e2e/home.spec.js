import { test, expect } from '@playwright/test';

test.describe('Главная страница после входа', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.fill('#login', 'testuser');
    await page.fill('#password', 'test123');
    await page.click('button:has-text("Войти")');
    await page.waitForURL('/', { timeout: 10000 });
  });

  test('Страница загружена', async ({ page }) => {
    // Проверяем, что есть какой-то контент
    const body = await page.textContent('body');
    expect(body).toBeTruthy();
    console.log('Страница загружена');
  });

  test('Есть кнопка "Выйти"', async ({ page }) => {
    await expect(page.locator('button:has-text("Выйти")')).toBeVisible({ timeout: 5000 });
  });
});
