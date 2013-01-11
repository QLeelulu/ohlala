define('#jquery/scroll-to-top', ['jquery'], function(require) {
    var $ = jQuery = require('jquery');

    /**
     * 返回顶部（代码来自极客公园）
     */
    var scroll_animate = upAnimate = toLeftFireAnimation = false, 
        anim_time = 500, rocketFireAnimateTime = 100, rocketFireFrameStart = 298,
        rocketFireFrameLength = 149, rocketFireState = [0, 0, 0, 1];
    var getScrollY = function (){var b=0,a=0;if(typeof(window.pageYOffset)=="number"){a=window.pageYOffset;b=window.pageXOffset}else{if(document.body&&(document.body.scrollLeft||document.body.scrollTop)){a=document.body.scrollTop;b=document.body.scrollLeft}else{if(document.documentElement&&(document.documentElement.scrollLeft||document.documentElement.scrollTop)){a=document.documentElement.scrollTop;b=document.documentElement.scrollLeft}}}return a};
    var resetScrollUpBtn = function (){$("#scrollTop .level-2").css({"background-position":"-149px 0px",display:"none"});$("#scrollTop").css({"margin-top":"-125px",display:"none"});upAnimate=false;clearTimeout(rocketFireTimer)};
    window.rocketFireAnimate = function () {
        for (i = 0; i < rocketFireState.length; i++) {
            if (rocketFireState[i] == 1) {
                rocketFireState[i] = 0;
                if (!toLeftFireAnimation) {
                    if ((i + 2) < rocketFireState.length) {
                        rocketFireState[i + 1] = 1
                    } else {
                        rocketFireState[0] = 1;
                        toLeftFireAnimation = true
                    }
                } else {
                    if ((i - 1) < 0) {
                        rocketFireState[1] = 1
                    } else {
                        rocketFireState[i - 1] = 1;
                        toLeftFireAnimation = false
                    }
                }
                break
            }
        }
        $("#scrollTop .level-2").css({
            "background-position": "-" + (rocketFireFrameStart + (i * rocketFireFrameLength)) + "px 0px",
            display: "block"
        });
        rocketFireTimer = setTimeout("rocketFireAnimate()", rocketFireAnimateTime)
    };
    var initScrollTop = function () {
        $("#scrollTop div.level-3").hover(function() {
            if ($.browser.msie) {
                this.parentNode.children[0].style.display = "block"
            } else {
                $(this.parentNode.children[0]).stop().fadeTo(500, 1)
            }
        },
        function() {
            if (upAnimate || scroll_animate) {
                return
            }
            if ($.browser.msie) {
                this.parentNode.children[0].style.display = "none"
            } else {
                $(this.parentNode.children[0]).stop().fadeTo(500, 0)
            }
        });
        $("#scrollTop div.level-3").click(function() {
            scroll_animate = true;
            $("#scrollTop .level-2").css({
                "background-position": "-298px 0",
                display: "block"
            });
            op = $.browser.opera ? $("html") : $("html, body");
            rocketFireTimer = setTimeout("rocketFireAnimate()", rocketFireAnimateTime);
            op.animate({
                scrollTop: 0
            },
            "slow",
            function() {
                scroll_animate = false;
                if (!upAnimate) {
                    upAnimate = true;
                    thisTop = $("#scrollTop")[0].offsetTop + 250;
                    $("#scrollTop").animate({
                        "margin-top": "-=" + thisTop + "px"
                    },
                    300,
                    function() {
                        resetScrollUpBtn()
                    })
                }
            })
        });
        window.onscroll = function() {
            if ((!scroll_animate) && (!upAnimate)) {
                body_elem = $("body")[0];
                window_elem = $("html")[0];
                if (window.innerHeight) {
                    wind_height = window.innerHeight;
                    wind_scroll = window.scrollY
                } else {
                    wind_height = document.documentElement.clientHeight;
                    wind_scroll = getScrollY()
                }
                elem = $("body")[0];
                scrollBtn = $("#scrollTop")[0];
                if ((elem) && (scrollBtn)) {
                    if ((scrollBtn.style.display == "none") && ((wind_height * 1.2) < wind_scroll)) {
                        scroll_animate = true;
                        $("#scrollTop").fadeIn(anim_time,
                        function() {
                            scroll_animate = false;
                            this.style.display = "block"
                        })
                    }
                    if ((scrollBtn.style.display == "block") && ((wind_height * 1.2) > wind_scroll)) {
                        scroll_animate = true;
                        $("#scrollTop").fadeOut(anim_time,
                        function() {
                            scroll_animate = false;
                            this.style.display = "none"
                        })
                    }
                }
            }
        };
    };
    initScrollTop();
});