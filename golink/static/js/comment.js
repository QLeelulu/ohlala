define(function(require, exports, module) {
    var $ = require('jquery');

    function commenting (t) {
        $('#comment-content, #new-comment, #r-comment-content, #btn-reply, #btn-cancel-reply').attr('disabled', true);
    }

    function commented (t) {
        $('#comment-content, #new-comment, #r-comment-content, #btn-reply, #btn-cancel-reply').removeAttr('disabled');
    }

    var bgFadeUp=function(element,red,green,blue){
        if(element.fade){
            window.clearTimeout(element.fade);
        }
        var cssValue = "rgb("+red+","+green+","+blue+")";
        $(element).css("background-color",cssValue);
        if(red == 255 && green == 255 && blue == 255){
            $(element).css("background-color","");
            return;
        }
        var newRed = red + Math.ceil((255-red)/10);
        var newGreen = green + Math.ceil((255-green)/10);
        var newBlue = blue + Math.ceil((255-blue)/10);
        var repeat = function(){
            bgFadeUp(element,newRed,newGreen,newBlue);
        };
        element.fade=window.setTimeout(repeat,100);
    };

    /**
     * 添加新评论
     */
    function initNewComment () {
        $('#new-comment').click(function (e) {
            var d = {
                'parent_id': 0,
                'link_id': $('#link-id').val(),
                'content': $.trim($('#comment-content').val())
            };
            $.ajax({
                type: 'post',
                dataType: 'json',
                data: d,
                url: '/link/' + d['link_id'] + '/ajax-comment',
                beforeSend: function(XMLHttpRequest){
                    commenting();
                },
                success: function(data, textStatus){
                    if (data && data.success) {
                        $('#comment-content').val('');
                        oh.Msg.success('评论成功');
                        var ele = $(data.commentHTML);
                        ele.prependTo('#comment-list');
                        bgFadeUp(ele[0], 255, 246, 1);
                        if ($('#no-comment-yet').length) {
                            $('#no-comment-yet').remove();
                        }
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            oh.Msg.error('评论失败: ' + data.errors);
                        }
                    } else {
                        oh.Msg.error('请求出错，请稍后重试');
                    }
                },
                complete: function(XMLHttpRequest, textStatus){
                    commented();
                },
                error: function(){
                    oh.Msg.error('请求出错，请稍后重试');
                }
            });
        });
    }

    /**
     * 回复评论
     */
    function initReplyComment () {
        $('#btn-cancel-reply').click(function () {
            $('#reply-form').hide();
        });

        $(document.body).on('click', '#comment-list .cm a.rp', function () {
        // $('#comment-list .cm a.rp').click(function () {
            var rid = $(this).closest('.cm').attr('data-id');
            if (!rid) { return }
            $('#reply-comment-id').val(rid);
            $('#reply-form').insertAfter($(this).parent()).show();
            $('#r-comment-content').focus();
        });

        $('#btn-reply').click(function (e) {
            var d = {
                'parent_id': $('#reply-comment-id').val(),
                'link_id': $('#link-id').val(),
                'content': $.trim($('#r-comment-content').val())
            };
            $.ajax({
                type: 'post',
                dataType: 'json',
                data: d,
                url: '/link/' + d['link_id'] + '/ajax-comment',
                beforeSend: function(XMLHttpRequest){
                    commenting();
                },
                success: function(data, textStatus){
                    if (data && data.success) {
                        $('#r-comment-content').val('');
                        $('#reply-form').hide();
                        oh.Msg.success('评论成功');

                        var ele = $(data.commentHTML);
                        ele.appendTo( $('#cm-' + d['parent_id']).find('>.ct') );
                        bgFadeUp(ele[0], 255, 246, 1);
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            oh.Msg.error('评论失败: ' + data.errors);
                        }
                    } else {
                        oh.Msg.error('请求出错，请稍后重试');
                    }
                },
                complete: function(XMLHttpRequest, textStatus){
                    commented();
                },
                error: function(){
                    oh.Msg.error('请求出错，请稍后重试');
                }
            });
        });
    }

    /**
     * 折叠评论
     */
    function initExpandComment () {
        $(document.body).on('click', '#comment-list .cm .ep', function (e) {
            var p = $(this).closest('.cm');
            if ($(this).attr('data-status') === '1') {
                $(this).text('[–]').removeAttr('data-status');
                p.removeClass('collapsed').find('> .vt, > .ct .tx, > .ct .ed, > .ct .cd').show();
                
            } else {
                $(this).attr('data-status', '1').text('[+]');
                p.addClass('collapsed').find('> .vt, > .ct .tx, > .ct .ed, > .ct .cd').hide();
            }
        });
    };

    /**
     * comment投票
     */
    function initVoteComment () {
        $(document.body).on('click', '#comment-list .cm .vt a', function () {
            var t = $(this), vt = 0;
            if (t.hasClass('up')) {
                vt = 1;
            } else if (t.hasClass('down')) {
                vt = 2;
            } else {
                return;
            }
            var cid = t.closest('.cm').attr('data-id');
            if (!cid) {return}
            $.ajax({
                url: '/vote/comment/' + cid + '/' + vt,
                type: "post",
                dataType: "json",
                beforeSend: function(xhr){
                    t.attr('disabled', true);
                },
                success: function(data, textStatus){
                    if (data && data.Success === true) {
                        var p = t.closest('.cm');
                        p.find('.vote a').removeClass('on');
                        t.addClass('on');
                        p.find('.ct .uif .v').text(data.VoteNum + '分')
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            oh.Msg.error( data.Errors ? data.Errors : '请求出错，请稍后重试');
                        }
                    } else {
                        oh.Msg.error('请求出错，请稍后重试');
                    }
                },
                complete: function(xhr, status){
                    t.removeAttr('disabled');
                },
                error: function(){
                    oh.Msg.error('请求出错，请稍后重试');
                }
            });
        });
    }
    /**
     * 追加评论
     */
	function initLoadMoreComment() {
		// $('#comment-list .ldmore a').unbind('click');
		$(document.body).on('click', '#comment-list .ldmore a', function () {

			var t = $(this)
			var pId = t.attr('pId');
		    var d = {
		        'except_ids': t.attr('exIds'),
		        'parent_path': t.attr('pp'),
		        'top_parent_id': t.attr('tId'),
		        'link_id': t.attr('lId'),
		        'sort_type': t.attr('srt')
		    };
			var oldValue = t.html();

			$.ajax({
		        url: '/comment/loadmore/',
		        type: "post",
		        dataType: "json",
		        data: d,
		        beforeSend: function(xhr){
		            t.html("<span style='color:red'>加载中...</span>");
		        },
		        success: function(data, textStatus){
		            if (data) {
						$("#comment-list div[lmid=lm" + pId + "]").remove();
						$("#comment-list div[pid=pid" + pId + "]").append(data.Html);

						// initLoadMoreComment();

		            } else {
		                oh.Msg.error('请求出错，请稍后重试11');
		            }
		        },
		        complete: function(xhr, status){
		            t.html(oldValue);
		        },
		        error: function(){
		            oh.Msg.error('请求出错，请稍后重试');
		        }
		    });
		});
	}

    exports.init = function() {
        initNewComment();
        initReplyComment();
        initExpandComment();
        initVoteComment();
		initLoadMoreComment();
    };
});


