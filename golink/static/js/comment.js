define(function(require, exports, module) {
    var $ = require('jquery');

    function commenting (t) {
        $('#comment-content, #new-comment, #r-comment-content, #btn-reply, #btn-cancel-reply').attr('disabled', true);
    }

    function commented (t) {
        $('#comment-content, #new-comment, #r-comment-content, #btn-reply, #btn-cancel-reply').removeAttr('disabled');
    }

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
                        alert('评论成功');
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            alert('评论失败: ' + data.errors);
                        }
                    } else {
                        alert('请求出错，请稍后重试');
                    }
                },
                complete: function(XMLHttpRequest, textStatus){
                    commented();
                },
                error: function(){
                    alert('请求出错，请稍后重试');
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

        $('#comment-list .cm a.rp').click(function () {
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
                        alert('评论成功');
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            alert('评论失败: ' + data.errors);
                        }
                    } else {
                        alert('请求出错，请稍后重试');
                    }
                },
                complete: function(XMLHttpRequest, textStatus){
                    commented();
                },
                error: function(){
                    alert('请求出错，请稍后重试');
                }
            });
        });
    }

    /**
     * 折叠评论
     */
    function initExpandComment () {
        $('#comment-list .cm .ep').click(function (e) {
            var p = $(this).closest('.cm');
            if ($(this).attr('data-status') === '1') {
                $(this).text('[–]').removeAttr('data-status');
                p.removeClass('collapsed').find('> .vt, > .ct .tx, > .ct .ed, > .ct .cd').show();
                
            } else {
                $(this).attr('data-status', '1').text('[+]');
                p.addClass('collapsed').find('> .vt, > .ct .tx, > .ct .ed, > .ct .cd').hide();
            }
        });
    }

    exports.init = function() {
        initNewComment();
        initReplyComment();
        initExpandComment();
    };
});