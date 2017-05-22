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
	jQuery('.stdtable button.toRelease').click(function(){
        var me = $(this)
        jConfirm("确定要发布上线？","发布上线", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/release/toRelease",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '上线失败');
                        } else {
                            jAlert('上线成功，请尽快验证', '上线成功', function(){
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
	jQuery('.stdtable button.releaseSuccess').click(function(){
        var me = $(this)
        jConfirm("确定要设置验证通过？验证通过的代码将合并到master并无法回滚","设置验证通过", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/release/releaseSuccess",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '设置验证通过失败');
                        } else {
                            jAlert('设置验证通过成功，代码已进入master分支', '设置验证通过成功', function(){
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
	jQuery('.stdtable button.releaseFail').click(function(){
        var me = $(this)
        jConfirm("确定要设置验证不通过？","设置验证不通过", function(c){
            if(c) {
                jPrompt("简短说明原因","","验证不通过",function(msg){
                    if (msg == null) {
                        return false
                    } else if (msg == "") {
                        jAlert('请输入原因', '验证不通过');
                        return false
                    } else {
                        $.ajax({
                            type: "POST",
                            dataType: "json",
                            url: "/release/releaseFail",
                            data: "id="+me.attr("value")+"&msg="+msg,
                            success: function(msg){
                                if (msg.err_no > 0)  {
                                    jAlert(msg.err_msg, '设置验证不通过失败');
                                } else {
                                    jAlert('设置验证不通过', '设置验证不通过', function(){
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
