from selenium import webdriver
from time import sleep
from dotenv import load_dotenv
import os 
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

    assert result == "sameorigin" or result == "deny" or result == "allow-from:"
    


