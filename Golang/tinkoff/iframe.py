from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options

options = Options()
options.add_argument("start-maximized")
options.add_argument("--no-sandbox")
options.add_argument("--disable-dev-shm-usage")
options.add_argument('--headless')
# создание экземпляра драйвера
driver = webdriver.Chrome()

# переход на страницу
driver.get("https://www.lamoda.ru/lp/tinkoff_loyalty_kredit/")

# ожидание появления iframe
iframe_locator = (By.XPATH, "//iframe[@id='tinkoff-iframe-form-wrapper']")
iframe = WebDriverWait(driver, 10).until(EC.presence_of_element_located(iframe_locator))

# переключение на iframe
driver.switch_to.frame(iframe)

# проверка наличия элементов в iframe
buttons_locator = (By.XPATH, "//button[contains(text(), 'Далее')]")
try:
    WebDriverWait(driver, 10).until(EC.visibility_of_element_located(buttons_locator))
except:
    print("Test FAILED: iframe not found")
    driver.quit()
    exit()

input_locator = (By.XPATH, "//input[@id='tveYLYapZp']")  # Ввод промокода
try:
    WebDriverWait(driver, 10).until(EC.visibility_of_element_located(input_locator))
except:
    print("Test FAILED: elements not found in iframe")
    driver.quit()
    exit()

# переключение на основное содержимое страницы
driver.switch_to.default_content()

# закрытие браузера
driver.quit()

print("Test OK")
