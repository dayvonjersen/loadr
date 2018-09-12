import codecs
import os
import sys
import PortableDevices as p

byteUnits=['','K','M','G','T','P']
def formatBytes(b):
    k=0
    while b>=1024:
        k+=1
        b/=1024
    return '%.1f %sB' % (b, byteUnits[k])

totalSize=int(sys.argv[1])
written=0
fileCount=0
def printProgress(currentFile, size, totalFiles):
    global written
    global fileCount
    written+=size
    fileCount+=1
    print('[%.1f%%] Wrote %d of %d files: %s / %s' % (float(written)/float(totalSize)*100.0, fileCount, totalFiles, formatBytes(written), formatBytes(totalSize)))

target="WALKMAN NWZ-E353/Storage Media/MUSIC"
content=p.getContentFromDevicePath(target)
with codecs.open("playlist.txt", "r", "utf-8") as f:
    playlist=f.readlines()
    for mp3 in playlist:
        mp3=mp3.strip()
        mp3Name=os.path.basename(mp3)
        mp3Size=os.path.getsize(mp3)
        mp3File=open(mp3, "rb")
        print('Copying %s (%s) ...' % (mp3Name.encode('utf-8'), formatBytes(mp3Size)))
        content.uploadStream(mp3Name, mp3File, mp3Size)
        mp3File.close()
        printProgress(mp3Name, mp3Size, len(playlist))

print('Done.')
