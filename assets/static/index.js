$(function() {

    var dl = $('dl.list');
    var dts = dl.find('dt');

    $('#search_field').on("keyup", function() {
        var box = $(this);
        var val = box.val();
        var reg = new RegExp(val, "gim");

        dts.each(function() {
            var dt = $(this);
            if (val == "") {
                dt.show();
            } else {
                var title = $("span.title", dt).html();
                var path  = $("a.link", dt).html();
                // console.log(title);
                // console.log(path)
                if (reg.test(title) || reg.test(path)) {
                    dt.show();
                } else {
                    dt.hide();
                }
            }
        });
    })

});
