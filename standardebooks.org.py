# m-we
# v1.0.1 (2018-11-18)
#
# Downloads books from standardebooks.org.

import os
import re
import urllib.request

# Returns a list of URLs to books listed on the standardebooks.org website.
def ScrapeEbooks():
    ebook_urls = []
    base_url = "https://standardebooks.org/ebooks/?page="
    i = 0
    flag = False
    while flag == False:
        i += 1
        flag = True
        url = base_url + str(i)
        print(url)
        html = urllib.request.urlopen(url).read().decode()
        for ebook_url in re.findall("/ebooks.*\">", html):
            ebook_url = ebook_url.replace("\">", "")
            if ebook_url != "/ebooks/" and ebook_url.find("?") == -1 and ebook_url.find("/ebooks/\"") == -1:
                ebook_url = "https://standardebooks.org" + ebook_url
                if not ebook_url in ebook_urls:
                    ebook_urls.append(ebook_url)
                    flag = False
    return ebook_urls

# Passed a list of URLs to standardebooks.org books, finds the download links
# for .epub, .epub3, and .azw3 files and downloads them.
def DownloadEbooks(ebook_urls, location):
    if not os.path.isdir(location):
        os.makedirs(location)
    for link in ebook_urls:
        html = urllib.request.urlopen(link).read().decode()
        for dl in re.findall("https://standardebooks.org/ebooks/.*(?=\")", html):
            if dl.endswith(".epub") or dl.endswith(".epub3") or \
               dl.endswith(".azw3"):
                print(dl)
                filename = dl[dl.rfind("/") + 1:]
                urllib.request.urlretrieve(dl, location + filename)

def main():
    DownloadEbooks(ScrapeEbooks(), "ebooks/")
    return 0

if __name__ == "__main__":
    main()
