define(function(require, exports, module) {
    var $ = require('jquery');

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

			$.ajax({
		        url: '/invite/email/',
		        type: "post",
		        dataType: "json",
		        data: {"emails":""},
		        beforeSend: function(xhr){
		            $('#fetchInviteUrl').attr('disabled', true);
		        },
		        success: function(data, textStatus){
		            if (data.Result == true) {
						$("#inviteUrl").val(data.InviteUrl);

						//$.copy(data.InviteUrl);
						oh.Msg.success('获取邀请url成功!');//,并复制到剪贴板
		            } else {
		                oh.Msg.error(data.Msg);
		            }
					$('#fetchInviteUrl').removeAttr('disabled');
		        },
		        complete: function(xhr, status){
					$('#fetchInviteUrl').removeAttr('disabled');
		        },
		        error: function(){
		            oh.Msg.error('请求出错，请稍后重试');
					$('#fetchInviteUrl').removeAttr('disabled');
		        }
		    });
		});
	}

    exports.init = function() {
        initSendInvite();
		initFetchInviteUrl();
    };
});


