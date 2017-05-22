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
	jQuery('.stdtable button.startTest').click(function(){
        var me = $(this)
        jConfirm("确定要部署测试？","部署测试", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/test/startTest",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '部署失败');
                        } else {
                            jAlert('部署成功，请尽快测试', '部署成功', function(){
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
	jQuery('.stdtable button.testSuccess').click(function(){
        var me = $(this)
        jConfirm("确定要设置测试通过？","设置测试通过", function(c){
            if(c) {
                $.ajax({
                    type: "POST",
                    dataType: "json",
                    url: "/test/testSuccess",
                    data: "id="+me.attr("value"),
                    success: function(msg){
                        if (msg.err_no > 0)  {
                            jAlert(msg.err_msg, '设置测试通过失败');
                        } else {
                            jAlert('设置测试通过成功，代码已进入待发布列表', '设置测试通过成功', function(){
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
	jQuery('.stdtable button.testFail').click(function(){
        var me = $(this)
        jConfirm("确定要设置测试不通过？","设置测试不通过", function(c){
            if(c) {
                jPrompt("简短说明原因","","测试不通过",function(msg){
                    if (msg == null) {
                        return false
                    } else if (msg == "") {
                        jAlert('请输入原因', '测试不通过');
                        return false
                    } else {
                        $.ajax({
                            type: "POST",
                            dataType: "json",
                            url: "/test/testFail",
                            data: "id="+me.attr("value")+"&msg="+msg,
                            success: function(msg){
                                if (msg.err_no > 0)  {
                                    jAlert(msg.err_msg, '设置测试不通过失败');
                                } else {
                                    jAlert('设置测试不通过', '设置测试不通过', function(){
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
