/*
 * 	Additional function for tables.html
 *	Written by ThemePixels	
 *	http://themepixels.com/
 *
 *	Copyright (c) 2012 ThemePixels (http://themepixels.com)
 *	
 *	Built for Amanda Premium Responsive Admin Template
 *  http://themeforest.net/category/site-templates/admin-templates
 */

jQuery(document).ready(function(){

	///// DELETE INDIVIDUAL ROW IN A TABLE /////
	jQuery('.stdtable a.delete').click(function(){
        var me = $(this)
        jConfirm("确定禁用？","禁用服务器", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/machine/del",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '禁用失败');
                        } else {
                            jAlert('禁用成功', '禁用成功', function(){
                                me.parent().prev().text("否")
                                me.parent().text("")
                            });
                        }
                    }
                });
            }
        })
		return false;
	});
	
});
