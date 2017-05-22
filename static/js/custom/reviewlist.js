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
	jQuery('.stdtable button.startReview').click(function(){
        var me = $(this)
        jConfirm("确定要获取Review权限？","获取Review权限", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/review/startReview",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '获取Review权限失败');
                        } else {
                            jAlert('获取Review权限成功，请尽快Review', '获取Review权限成功', function(){
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
	jQuery('.stdtable button.reviewSuccess').click(function(){
        var me = $(this)
        jConfirm("确定要设置Review通过？","设置Review通过", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/review/reviewSuccess",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '设置Review通过失败');
                        } else {
                            jAlert('设置Review通过成功', '设置Review通过成功', function(){
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
	jQuery('.stdtable button.reviewFail').click(function(){
        var me = $(this)
        jConfirm("确定要设置Review不通过？","Review不通过", function(c){
            if(c) {
                jPrompt("简短说明原因","","Review不通过",function(msg){
                    if (msg == null) {
                        return false
                    } else if (msg == "") {
                        jAlert('请输入原因', 'Review不通过');
                        return false
                    } else {
                        $.ajax({
                            type: "POST",
                            dataType: "json",
                            url: "/review/reviewFail",
                            data: "id="+me.attr("value")+"&msg="+msg,
                            success: function(msg){
                                if (msg.err_no > 0)  {
                                    jAlert(msg.err_msg, 'Review不通过');
                                } else {
                                    jAlert('设置Review不通过', '设置Review不通过', function(){
                                        window.location.reload()
                                    });
                                }
                            }
                        });
                    }
                })
            }
        })
        return false;
	});
	
	///// TRANSFORM CHECKBOX AND RADIO BOX USING UNIFORM PLUGIN /////
	jQuery('input:checkbox,input:radio').uniform();
	
	
});
