(function (w, d, id) {
    var gbk2utf8 = function(gbk){
        if(!gbk){return '';}
        var utf8 = [];
        for(var i=0;i<gbk.length;i++){
            var s_str = gbk.charAt(i);
            if(!(/^%u/i.test(escape(s_str)))){utf8.push(s_str);continue;}
            var s_char = gbk.charCodeAt(i);
            var b_char = s_char.toString(2).split('');
            var c_char = (b_char.length==15)?[0].concat(b_char):b_char;
            var a_b =[];
            a_b[0] = '1110'+c_char.splice(0,4).join('');
            a_b[1] = '10'+c_char.splice(0,6).join('');
            a_b[2] = '10'+c_char.splice(0,6).join('');
            for(var n=0;n<a_b.length;n++){
                utf8.push('%'+parseInt(a_b[n],2).toString(16).toUpperCase());
            }
        }
        return utf8.join('');
    };
    function openSubmit () {
        var t = d.title, u=d.location.href;
        t = t.split(' - ')[0];
        w.open('http://127.0.0.1:8080/link/submit?url='+u+'&title='+gbk2utf8(t));
    }
    w[id + '_fun__submit'] = openSubmit;
    // 这样打开，弹出窗会被拦，闷
    openSubmit();
    var j = d.getElementById(id + '__script__tag');
    j.parentNode.removeChild(j);
})(window, document, '__ohlala__')
/*
function g2u(gbk){if(!gbk){return'';}
var utf8=[];for(var i=0;i<gbk.length;i++){var s_str=gbk.charAt(i);if(!(/^%u/i.test(escape(s_str)))){utf8.push(s_str);continue;}var s_char=gbk.charCodeAt(i);var b_char=s_char.toString(2).split('');var c_char=(b_char.length==15)?[0].concat(b_char):b_char;var a_b=[];a_b[0]='1110'+c_char.splice(0,4).join('');a_b[1]='10'+c_char.splice(0,6).join('');a_b[2]='10'+c_char.splice(0,6).join('');for(var n=0;n<a_b.length;n++){utf8.push('%'+parseInt(a_b[n],2).toString(16).toUpperCase());}}return utf8.join('');};
 */