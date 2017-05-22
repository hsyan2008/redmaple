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
	jQuery('.stdtable button.toReview').click(function(){
        var me = $(this)
        jConfirm("确定提交Review？","提交Review", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/task/toReview",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '提交失败');
                        } else {
                            jAlert('提交成功', '提交成功', function(){
                                window.location.reload()
                            });
                        }
                    }
                });
            }
        })
        return false;
	});
	///// DELETE INDIVIDUAL ROW IN A TABLE /////
	jQuery('.stdtable button.delete').click(function(){
        var me = $(this)
        jConfirm("确定删除？将清理本次任务的所有分支代码","删除任务", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/task/del",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '删除失败');
                        } else {
                            jAlert('删除成功', '删除成功', function(){
                                window.location.reload()
                            });
                        }
                    }
                });
            }
        })
        return false;
	});
	
	///// TRANSFORM CHECKBOX AND RADIO BOX USING UNIFORM PLUGIN /////
	jQuery('input:checkbox,input:radio').uniform();
	
	
});
