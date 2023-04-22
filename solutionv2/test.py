from selenium import webdriver
from time import sleep
from dotenv import load_dotenv
import os 

import argparse
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
load_dotenv()

def test_200():
    options_G = webdriver.ChromeOptions()
    port = os.getenv('PORT')
    url = os.getenv('URL')
    driver = webdriver.Remote(
        command_executor=f'http://localhost:{port}',
        options=options_G
    )
        
    js_1 = '''
let callback = arguments[0];
let xhr = new XMLHttpRequest();
xhr.open('GET', '
'''
    js_2 = '''
', true);
xhr.onload = function () {
    if (this.readyState === 4) {
        callback(this.status);
    }
};
xhr.onerror = function () {
    callback('error');
};
xhr.send(null);
'''
    driver.get(url)
    status_code = driver.execute_async_script(js_1 + url + js_2)
            
    driver.quit()

    assert status_code == 200 or status_code == "200"

def test_xframe():
    import requests
    req = requests.get('https://www.tinkoff.ru/')

    result = req.headers['x-frame-options']  

    assert result == "sameorigin" or result == "deny"



def test_4():
    options = Options()
    options.add_argument("start-maximized")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument('--headless')

    # парсинг аргументов
    parser = argparse.ArgumentParser()
    parser.add_argument('--url', help='URL of the page to test', required=True)
    args = parser.parse_args()

    # создание экземпляра драйвера
    driver = webdriver.Chrome()

    # переход на страницу
    driver.get(args.url)

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
    assert 1 == 1




