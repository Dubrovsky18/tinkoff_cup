from selenium import webdriver
# from selenium.webdriver.common.keys import Keys
from time import sleep
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager


options = Options()
options.add_argument("start-maximized")
options.add_argument("--no-sandbox")
options.add_argument("--disable-dev-shm-usage")
options.add_argument('--headless')
driver = webdriver.Chrome(service=Service("./chromedriver"), options=options)
driver.get("https://translate.google.com/?hl=ru&sl=ru&tl=en&op=translate")
sleep(2)
#
textarea = driver.find_element(by=By.XPATH, value='//*[@id="yDmH0d"]/c-wiz/div/div[2]/c-wiz/div[2]/c-wiz/div[1]/div[2]/div[3]/c-wiz[1]/span/span/div/textarea')
textarea.click()
sleep(1)
textarea.send_keys('Тинькофф Кубок')
sleep(3)
translated_text = driver.find_element(by=By.XPATH, value='//*[@id="yDmH0d"]/c-wiz/div/div[2]/c-wiz/div[2]/c-wiz/div[1]/div[2]/div[3]/c-wiz[2]/div/div[9]/div/div[1]/span[1]/span/span')
result = translated_text.text
if result == "Tinkoff Cup":
    print("Test OK")
else:
    print("Test FAILED")
driver.close()