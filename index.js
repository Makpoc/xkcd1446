(function() {
    let IMGPREFIX= 'data:image/png;base64,'

    try {
        let source = new EventSource('/landing');
        source.addEventListener('landing', function(ev){
            let newimage = IMGPREFIX + ev.data;
            let image = document.getElementById('landing').getElementsByTagName('img')[0];
            if (image === undefined ) {
                image = document.createElement('img');
                document.getElementById('landing').appendChild(image);
            }
            image.src = newimage;
        }, false);
    } catch (e) {
        console.log('js_error', e)
    }
})()
