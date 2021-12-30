# convert -size 360x360 xc:white -font "FreeMono" -pointsize 12 -fill black -draw @ascii.txt image.png
convert -size 596x620 xc:black -font "Monaco" -pointsize 12 -fill green -annotate +15+15 "@001-final.txt" 001-final.txt.png
magick Micromove-nebula-16.jpg foo.png -gravity center -compose over -composite resultimage.png
magick Micromove-nebula-16.jpg foo.png -gravity center -compose over -composite resultimage.png


find convert dragon_sm.gif    -resize 64x64  resize_dragon.gif

convert -delay 5 -loop 0 *.png 001.gif

find . -type f -name '*.tif' | xargs -I xxx convert xxx -resize 1000x620 xxx_sm.jpg



find out -type f -name '*.txt' | awk '{print "convert -size 596x620 xc:black -font \"DejaVu-Sans-Mono\" -pointsize 12 -fill green -annotate +15+15 \"@" $1 "\" " $1 ".png"'} | parallel
find out -type f -name '*-final.txt' | awk '{print "convert -size 596x620 xc:black -font \"DejaVu-Sans-Mono\" -pointsize 12 -fill transparent -annotate +15+15 \"@" $1 "\" " $1 "-alpha.png"'} | parallel

convert mask.jpg -gravity center -crop 596x620  mask-crop.jpg
convert mask-crop-0.jpg 000-final.txt-alpha.png -gravity center -size 596x620 -compose over -composite maze-color.png