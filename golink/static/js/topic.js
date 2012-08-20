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
        $('#fileupload').fileupload({
            dataType: 'json',
            paramName: 'topic-image',
            done: function (e, data) {
                $.each(data.result, function (index, file) {
                    $('<p/>').text(file.name).appendTo(document.body);
                });
            }
        });
    });

});