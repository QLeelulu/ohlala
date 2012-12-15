define(function(require, exports, module) {
    var $ = require('jquery');
    require('jquery.fileupload');

    // $(".dofollow").click(function () {
    //     var btn = $(this);
    //     btn.attr("disabled", true).text("关注中...")
    //     $.ajax({
    //         url: "/topic/" + btn.attr("data-uid") + "/follow",
    //         type: "post",
    //         dataType: "json",
    //         success: function (r) {
    //             if (r && r.success){
    //                 btn.text("已关注");
    //             }else {
    //                 btn.text("关注");
    //                 btn.removeAttr("disabled");
    //                 alert(r.errors)
    //             }
    //         }
    //     });
    // });

    $(function () {
        $('.topic-img-upload').fileupload({
            dataType: 'json',
            paramName: 'topic-image',
            done: function (e, data) {
                if (data.result) {
                    if (data.result.success) {
                        oh.Msg.success("上传成功");
                    } else {
                        oh.Msg.error('上传出错: ' + data.result.errors);
                    }
                } else {
                    oh.Msg.error('上传出错');
                }
            }
        });
    });

});