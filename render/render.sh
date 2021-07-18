# Create the full tree
go run render.go

# crop the masks
find out -mindepth 1 -type d | awk '{print "convert " $1 "/mask.jpg -gravity center -crop 596x620  " $1 "/mask-crop.jpg"'} | parallel

# Create every frame
find out -type f -name '*.txt'       | awk '{print "convert -size 596x620 xc:black -font \"DejaVu-Sans-Mono\" -pointsize 12 -fill transparent -annotate +15+15 \"@" $1 "\" " $1 ".png"'} | parallel

# Create color frames
find out -mindepth 2 -type f -name '*.txt.png' | sed 's/.png//' | sort | awk -F '/' '{print "convert " $1 "/" $2 "/mask-crop-0.jpg " $1 "/" $2 "/" $3 ".png -gravity center -size 596x620 -compose over -composite " $1 "/" $2 "/" $3 "-color.png"}' | parallel

# render gifs
# skip this for now since we have MP4s
# find out -mindepth 1 -type d | awk '{print "convert -delay 5 -loop 0 " $1 "/*.txt-color.png " $1 "/animated.gif"}' | parallel

# render mp4s
find out -mindepth 1 -type d | awk '{print "ffmpeg -framerate 30 -pattern_type glob -i \"" $1 "/*.txt-color.png\" -c:v libx264 -r 30 -pix_fmt yuv420p " $1 "/maze.mp4"}' | parallel

# final cleanup
Xvfb :1 -screen 0 1920x1080x24 &
export DISPLAY=:1

while read -u 3 line; do
    pushd $line
    mkdir final
    mv *-final.txt-color.png final/maze-color.png
    cp final/maze-color.png final/maze.png
#   mv animated.gif final/maze.gif
    mv *.json final/
    mv *.sha1 final/maze.sha1
    cp *-final.txt.png final/
    mv *-final.txt.png final/maze-plain.png
    mv *-final.txt final/maze.txt
    mv maze.scad final/maze.scad
    mv maze.mp4 final/maze.mp4
    mv index.html final/index.html

    find . -maxdepth 1 -type f | xargs rm
    mv final/* .
    openscad -o maze-3d.png maze.scad --camera=19.34,24.2,3.47,25.2,0,135.1,147.57 --projection=ortho --imgsize=3840,2160 --colorscheme=Starnight
#   openscad -o maze.stl maze.scad

    rm -rf final/
    popd
done 3<<< $(find out -mindepth 1 -maxdepth 1 -type d)

find out -mindepth 1 -type d | awk '{print "openscad -o " $1 "/maze.stl " $1 "/maze.scad"}' | parallel
