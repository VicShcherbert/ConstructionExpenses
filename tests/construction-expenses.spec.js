import { test, expect } from '@playwright/test';

test('title-test', async ({ page }) => {
  await page.goto('localhost:3000');
  await expect(page).toHaveTitle('Constructifyer Projects');
});

test('click-login-test', async ({page}) => {
    await page.goto('localhost:3000');
    await page.getByRole('link', {name: 'Login'}).click();
})

test('click-fairview-project', async ({page}) => {
    await page.goto('localhost:3000/projects');
    await page.getByRole('button', {name: '4017 E. Fairview Ave.'}).click();
    const locator = page.locator('.title')
    await expect(locator).toContainText('4017 E. Fairview Ave.');
    await expect(locator).toContainText('Total Expenses');
})
