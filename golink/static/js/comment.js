define(function(require, exports, module) {
    var $ = require('jquery');

    function showCommenting () {
        $('#comment-content, #new-comment').attr('disabled', true);
    }

    function hideCommenting () {
        $('#comment-content, #new-comment').removeAttr('disabled');
    }

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
                    showCommenting();
                },
                success: function(data, textStatus){
                    if (data && data.success) {
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
                    hideCommenting();
                },
                error: function(){
                    alert('请求出错，请稍后重试');
                }
            });
        });
    }

    exports.init = function() {
        initNewComment();
    };
});