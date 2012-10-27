define(function(require, exports, module) {
    var $ = require('jquery');

    function inviting (t) {
        $('#comment-content, #new-comment, #r-comment-content, #btn-reply, #btn-cancel-reply').attr('disabled', true);
    }

    function invited (t) {
        $('#comment-content, #new-comment, #r-comment-content, #btn-reply, #btn-cancel-reply').removeAttr('disabled');
    }

    /**
     * 发送邀请
     */
	function initSendInvite() {

		$(document.body).on('click', '#sendInvite', function () {

			var emails = $.trim($("#inviteEmails").val());
			if (emails == "") {
				oh.Msg.error('请输入邀请的Email地址!');
			}
			var arrEmails = emails.split(/；|;/);
			var reg = /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$/;
			var arrLegalEmails = [];

			for(i in arrEmails){
				var email = $.trim(arrEmails[i])
				if (email != "" && reg.test(email) == false) {
					oh.Msg.error('输入的Email地址有误,请重输!');
					return;
				}
				if (email != "") {
					arrLegalEmails.push(email);
				}
			}

			$.ajax({
		        url: '/invite/email/',
		        type: "post",
		        dataType: "json",
		        data: {"emails":arrLegalEmails.join(';')},
		        beforeSend: function(xhr){
		            $('#sendInvite').attr('disabled', true);
		        },
		        success: function(data, textStatus){
		            if (data.Result == true) {
						$("#inviteEmails").val("");
						oh.Msg.success('邀请成功');

		            } else {
		                oh.Msg.error(data.Msg);
		            }
					$('#sendInvite').removeAttr('disabled');
		        },
		        complete: function(xhr, status){
					$('#sendInvite').removeAttr('disabled');
		        },
		        error: function(){
		            oh.Msg.error('请求出错，请稍后重试');
					$('#sendInvite').removeAttr('disabled');
		        }
		    });
		});
	}

    /**
     * 获取邀请url
     */
	function initFetchInviteUrl() {

		$(document.body).on('click', '#fetchInviteUrl', function () {

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
        initSendInvite();
		initFetchInviteUrl();
    };
});


