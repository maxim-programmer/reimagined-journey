import { test, expect } from '@playwright/test';

test.describe('Простая проверка', () => {
  test('Страница загружается', async ({ page }) => {
    await page.goto('http://localhost');
    // Проверяем, что страница загрузилась - любой текст
    const body = await page.textContent('body');
    expect(body).toBeTruthy();
    console.log('Страница загружена, содержимое:', body?.substring(0, 200));
  });

  test('На странице есть кнопка "Войти"', async ({ page }) => {
    await page.goto('http://localhost');
    // Ищем кнопку с текстом "Войти" любым способом
    const button = await page.locator('button:has-text("Войти")');
    await expect(button).toBeVisible({ timeout: 5000 });
    console.log('Кнопка "Войти" найдена');
  });
});
