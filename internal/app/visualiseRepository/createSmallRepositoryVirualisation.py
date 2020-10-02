# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

# http://cloc.sourceforge.net/

import subprocess


def main():
    visualisationVideoName = "repositoryVisualisation"
    ppmFile = f'{visualisationVideoName}.ppm'
    wmvFile = f'{visualisationVideoName}.wmv'

    ffmpegArray = ['ffmpeg','-y','-r','60','-f','image2pipe','-vcodec','ppm','-i',ppmFile,'-vcodec','wmv1','-r','60','-qscale','0', wmvFile]  # https://www.ffmpeg.org/download.html#build-windows
    subprocess.run(ffmpegArray)


if __name__ == '__main__':
    main()