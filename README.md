# GetArtWork
Get artwork image from Museum

Now support:
* MET(Metropolitan Museum)

# Usage
GetArtWork.exe -s=MET -a="van gogh" -t=sunflower -d -m=Painting -p=D:/images
* s: source, from which museum
* a: artist name
* t: artwork title
* d: download or not, if not set, images will not download, but only preview what you will download
* m: medium, now support Painting, Drawing, Print, Album, Tapestry, Printtest
* p: download path, if download, will create subdirctory by artist name
