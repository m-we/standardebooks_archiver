import os
import re
import sys
import urllib.request

def find():
    base_url = 'https://standardebooks.org/ebooks?page='
    urls = set()
    page = 1
    while True:
        print('Checking page ' + str(page))
        urls_len = len(urls)
        resp = urllib.request.urlopen(base_url + str(page)).read().decode()
        for url in re.findall('/ebooks/.*/.+?">', resp):
            url = url.replace('">', '').replace('" tabindex="-1','')
            urls.add('https://standardebooks.org' + url)
        if len(urls) == urls_len:
            print('\tdone\n')
            return urls
       page += 1

def download(urls, frmt, folder):
    if not os.path.isdir(folder):
        os.makedirs(folder)
    for url in urls:
        resp = urllib.request.urlopen(url).read().decode().replace('\\','')
        for dlink in re.findall('https://standardebooks.org/ebooks/.*(?=")', resp):
            if dlink.endswith(frmt):
                if (frmt == 'epub' or frmt == '.epub') and (dlink.find('kepub') != -1 or dlink.find('_advanced') != -1):
                    continue
                
                print(dlink)
                fn = folder + os.path.basename(dlink)
                if not os.path.isfile(fn):
                    urllib.request.urlretrieve(dlink, fn)

if __name__ == '__main__':
    u = find()
    download(u, sys.argv[1], sys.argv[2] + '/')

