<template>
<div>
    <!-- 添加项目 -->
    <div id="add_project">
        <el-button type="primary" @click="new_product">新建项目</el-button>
        <el-input placeholder="项目名称" v-model="proname" style="width:200px"></el-input>
        <el-input placeholder="项目描述" v-model="prodes" style="width:200px"></el-input>
        <br><br>
    </div>
    <!-- 项目列表 -->
    <el-table
        stripe
        border
        :data="projects.filter(data => !search || data.proname.toLowerCase().includes(search.toLowerCase()) || data.prodes.toLowerCase().includes(search.toLowerCase()))"
        empty-text="暂无项目"
        :header-cell-style="{background:'#ddd'}"
        highlight-current-row>
        <el-table-column label="项目名称" align="center" prop="proname" sortable>
        </el-table-column>
        <el-table-column label="项目描述" align="center" prop="prodes">
        </el-table-column>
        <el-table-column label="最近修改" align="center" prop="update_time" sortable>
            <template slot-scope="scope">
                <p>{{scope.row.update_time|dateFormat}}</p>
            </template>
        </el-table-column>
        <el-table-column label="前端" align="center" width="100">
            <template slot-scope="scope">
                <el-tooltip class="item" effect="dark" content="点击加入前端测试" placement="top">
                    <a href="#" @click.prevent="go_webTest(scope.row)">
                        <p v-if="scope.row.webresult" style="color:green">PASS</p>
                        <p v-else style="color:red">FAIL</p>
                    </a>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column label="后端" align="center" width="100">
            <template slot-scope="scope">
                <el-tooltip class="item" effect="dark" content="点击加入后端测试" placement="top">
                    <a href="#" @click.prevent="go_apiTest(scope.row)">
                        <p v-if="scope.row.apiresult" style="color:green">PASS</p>
                        <p v-else style="color:red">FAIL</p>
                    </a>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column label="移动端" align="center" width="100">
            <template slot-scope="scope">
                <el-tooltip class="item" effect="dark" content="点击加入移动端测试" placement="top">
                    <a href="#" @click.prevent="go_appTest(scope.row)">
                        <p v-if="scope.row.appresult" style="color:green">PASS</p>
                        <p v-else style="color:red">FAIL</p>
                    </a>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column label="测试结果" align="center" width="100">
            <template slot-scope="scope">
                <el-tooltip class="item" effect="dark" content="点击开始项目测试(前端+后端)" placement="top">
                    <a href="#" @click.prevent="projectTest(scope.row)">
                        <p v-if="scope.row.result" style="color:green" v-html="passText"></p>
                        <p v-else style="color:red" v-html="failText"></p>
                    </a>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column align="center" width="200px">
            <template slot="header">
                <el-input v-model="search" size="mini" placeholder="输入关键字搜索"/>
            </template>
            <template slot-scope="scope">
                <el-tooltip class="item" effect="dark" content="编辑修改" placement="top">
                    <el-button
                        v-if='scope.row.user==userId || userId==1'
                        size="mini"
                        type="primary"
                        @click="open_edit(scope.row)" class="el-icon-edit">
                    </el-button>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="测试报告" placement="top">
                    <el-button
                        size="mini"
                        type="primary"
                        @click="go_report(scope.row)" icon="el-icon-document">
                    </el-button>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="删除" placement="top">
                    <el-button
                        v-if='scope.row.user==userId || userId==1'
                        size="mini"
                        type="danger"
                        @click="handleDelete(scope.$index, scope.row)" icon="el-icon-delete">
                    </el-button>
                </el-tooltip>
            </template>
        </el-table-column>
    </el-table>
    <br>
    <!-- 翻页 -->
    <div style="text-align: center;">
        <el-button type="primary" :disabled="isPreDisabled" @click="get_pre">上一页</el-button>
        <el-button type="primary" :disabled="isNextDisabled" @click="get_next">下一页</el-button>
    </div>
    <!-- 修改数据 -->
    <el-dialog :visible.sync="dialogFormVisible">
        <el-form>
            <el-form-item label="项目名称" label-width="120px">
            <el-input v-model="editObj.proname" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <el-form >
            <el-form-item label="项目描述" label-width="120px">
            <el-input v-model="editObj.prodes" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="dialogFormVisible = false">取 消</el-button>
            <el-button type="primary" @click="handleEdit(editObj)">确 定</el-button>
        </div>
    </el-dialog>
</div>
</template>

<script>
/* eslint-disable */
import md5 from 'js-md5'
export default {
	name: "Project",
	data() {
		return {
            axios: this.axios,
            url: this.url,
            userId: this.storage.getItem('userID'),
            token: this.storage.getItem('token'),
            testUserId:this.storage.getItem('testUserId'),
            testUserToken:this.storage.getItem('testUserToken'),
            adminUserId:this.storage.getItem('adminUserId'),
            adminUserToken:this.storage.getItem('adminUserToken'),
            adminUser:'',
            adminPsw:'',
            testUser:'',
            testPsw:'',
            testBtn:'开始测试',
            loading:false,
            passText:'PASS',
            failText:'FAIL',
            proname: '',
            prodes: '',
            projects:[],
            search: '',
            pre:'',
            next:'',
            isNextDisabled:false,
            isPreDisabled:false,
            dialogFormVisible:false,
            editObj:{
                id:'',
                proname:'',
                prodes:'',
            }
		};
	},
	methods: {
        // 项目测试
        projectTest(row) {
            this.$confirm('项目测试需要等待时间较长, 是否开始测试?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
                }).then(() => {
                    this.$message({
                        message: '测试开始',
                        type: 'success',
                        center: true,
                        showClose: true
                    });
                    this.passText = this.failText = '<p class="el-icon-loading"></p>'
                    var params_data = {
                        'userID':this.userId,
                        'token':this.token,
                        'projectId':row['id']
                    }
                    var data = {}
                    if (row.proname == '量产云平台') {
                        var data = {
                            testUserInfo:{
                                testUserId:this.testUserId,
                                testUserToken:this.testUserToken,
                                adminUserId:this.adminUserId,
                                adminUserToken:this.adminUserToken,
                            }
                        }
                    }
                    this.axios({
                        baseURL:this.url,
                        url:'/api/v1/projectTest/',
                        method:'post',
                        params:params_data,
                        data:data
                    }).then(response=>{
                        // 判断是否成功
                        if (!response.data.errcode) {
                            this.$message({
                                message: row['proname'] + ' 测试 PASS',
                                type: 'success',
                                center: true,
                                showClose: true,
                                duration:0,
                            });
                        }
                        else {
                            this.$message({
                                message: response.data.errmsg,
                                type: 'error',
                                center: true,
                                showClose: true,
                                duration:0,
                            })
                        }
                        this.passText='PASS'
                        this.failText='FAIL'
                        this.get_projects()
                    },error=>{
                        this.$message({
                            message: '自动化测试平台异常，请检查网络',
                            type: 'error',
                            center: true,
                            showClose: true,
                            duration:0,
                        })
                        this.passText='PASS'
                        this.failText='FAIL'
                        this.testBtn='开始测试'
                    })
                }).catch(() => {
            });
        },
        // 加载数据
        get_projects() {
            var params_data = {'userID':this.userId,'token':this.token}
            this.axios({
                baseURL:this.url,
                url:'api/v1/project',
                method:'get',
                params:params_data,
            }).then(response=>{
                this.projects=response.data.datas
            },error=>{
                this.$message({
                        message: '匿名用户，请先登录',
                        type: 'error',
                        center: true,
                        showClose: true,
                    })
                this.$router.push('/')
            })
        },
        // 打开编辑
        open_edit(row) {
            this.dialogFormVisible = true
            // this.editObj=row
            this.editObj['id']=row.id
            this.editObj['proname']=row.proname
            this.editObj['prodes']=row.prodes
        },
        // 编辑修改数据
        handleEdit(row) {
            if (!row.proname || !row.prodes) {
                this.$message({
                    message: "项目名称和描述不能为空",
                    type: 'error',
                    center: true
                })
                return
            }
            // 关闭编辑框
            this.dialogFormVisible = false
            var params_data = {'userID':this.userId,'token':this.token}
            this.axios({
                baseURL:this.url,
                url:'api/v1/project/'+row.id,
                method:'patch',
                params:params_data,
                data:row,
            }).then(response=>{
                // 判断是否成功
                if (!response.data.errcode) {
                    this.$message({
                        message: '修改成功',
                        type: 'success',
                        center: true
                    });
                    this.get_projects()
                }
                else {
                    this.$message({
                        message: "修改失败",
                        type: 'error',
                        center: true
                    })
                }
            },error=>{
                this.$message({
                    message: '自动化测试平台异常，请检查网络',
                    type: 'error',
                    center: true
                })
            })
        },
        // 删除数据
        handleDelete(index, row) {
            this.$confirm('此操作将永久删除该项, 是否继续?', '提示', {
                distinguishCancelAndClose: true,
                type: 'warning',
                confirmButtonText: '确定',
                cancelButtonText: '取消'
            }).then(() => {
                var params_data = {'userID':this.userId,'token':this.token}
                this.axios({
                    baseURL:this.url,
                    url:'api/v1/project/'+row.id,
                    method:'delete',
                    params:params_data,
                }).then(response=>{
                    // 判断是否成功
                    if (!response.data.errcode) {
                        this.$message({
                            message: '删除成功',
                            type: 'success',
                            center: true
                        });
                    }
                    else {
                        this.$message({
                            message: "删除失败",
                            type: 'error',
                            center: true
                        })
                    }
                    this.get_projects()
                },error=>{
                    this.$message({
                        message: '自动化测试平台异常，请检查网络',
                        type: 'error',
                        center: true
                    })
                })
            }).catch(action=>{
                // console.log(action)
            })
        },
        // 添加数据
        new_product() {
            if (!this.proname || !this.prodes) {
                this.$message({
                    message: "项目名称和描述不能为空",
                    type: 'error',
                    center: true
                })
                return
            }
            var body_data = {
                    'proname': this.proname,
                    'prodes': this.prodes,
                    'user': parseInt(this.userId),
                }
            var params_data = {'userID':this.userId,'token':this.token}
            this.axios({
                baseURL:this.url,
                url:'api/v1/project',
                method:'post',
                params:params_data,
                data:body_data,
            }).then(response=>{
                // 判断是否成功
                if (!response.data.errcode) {
                    this.$message({
                        message: '新建成功',
                        type: 'success',
                        center: true
                    });
                // 重新加载数据
                this.get_projects()
                }
                else {
                    this.$message({
                        message: "新建失败",
                        type: 'error',
                        center: true
                    })
                }
            },error=>{
                this.$message({
                    message: '自动化测试平台异常，请检查网络',
                    type: 'error',
                    center: true
                })
            })
            this.proname=''
            this.prodes=''
        },
        // 上一页
        get_pre() {
            this.axios({
                url:this.pre,
                method:'get',
            }).then(response=>{
                // 判断是否成功
                if (!response.data.errcode) {
                    this.projects=response.data.results
                    // 判断是否有上一页
                    this.pre=response.data.previous
                    if (!this.pre) {
                        this.isPreDisabled=true
                    }
                    else {
                        this.isPreDisabled=false
                    }
                    // 判断是否有下一页
                    this.next=response.data.next
                    if (!this.next) {
                        this.isNextDisabled=true
                    }
                    else {
                        this.isNextDisabled=false
                    }
                }
                else {
                    this.$message({
                        message: "加载失败",
                        type: 'error',
                        center: true,
                        showClose: true,
                    })
                }
            },error=>{
                this.$message({
                        message: '匿名用户，请先登录',
                        type: 'error',
                        center: true,
                        showClose: true,
                    })
                this.$router.push('/')
            })
        },
        // 下一页
        get_next() {
            this.axios({
                url:this.next,
                method:'get',
            }).then(response=>{
                // 判断是否成功
                if (!response.data.errcode) {
                    this.projects=response.data.results
                    // 判断是否有上一页
                    this.pre=response.data.previous
                    if (!this.pre) {
                        this.isPreDisabled=true
                    }
                    else {
                        this.isPreDisabled=false
                    }
                    // 判断是否有下一页
                    this.next=response.data.next
                    if (!this.next) {
                        this.isNextDisabled=true
                    }
                    else {
                        this.isNextDisabled=false
                    }
                }
                else {
                    this.$message({
                        message: "加载失败",
                        type: 'error',
                        center: true,
                        showClose: true,
                    })
                }
            },error=>{
                this.$message({
                        message: '匿名用户，请先登录',
                        type: 'error',
                        center: true,
                        showClose: true,
                    })
                this.$router.push('/')
            })
        },
        // 前端测试
        go_webTest(project) {
            var url = '/home/test/'
            var query = {
                'projectId': project.id,
                'projectName':project.proname,
                'srcType':'web',
                'type':'前端'
            }
            this.$router.push({ path: url,query:query})
        },
        // 后端测试
        go_apiTest(project) {
            var url = '/home/test/'
            var query = {
                'projectId': project.id,
                'projectName':project.proname,
                'srcType':'api',
                'type':'后端'
            }
            this.$router.push({ path: url, query:query})
        },
        // 移动端测试
        go_appTest(project) {
            var url = '/home/test/'
            var query = {
                'projectId': project.id,
                'projectName':project.proname,
                'srcType':'app',
                'type':'移动端'
            }
            this.$router.push({ path: url, query:query})
        },
        // 测试报告
        go_report(row) {
            var url = '/home/report/'
            var query = {
                projectName:row.proname
            }
            this.$router.push({ path: url, query:query})
        }
    },
    beforeCreate() {
    },
    created() {
        this.get_projects()
    },
    filters:{
        dateFormat:function(time) {
            var date=new Date(time);
            var year=date.getFullYear();
            /* 在日期格式中，月份是从0开始的，因此要加0
            * 使用三元表达式在小于10的前面加0，以达到格式统一  如 09:11:05
            * */
            var month= date.getMonth()+1<10 ? "0"+(date.getMonth()+1) : date.getMonth()+1;
            var day=date.getDate()<10 ? "0"+date.getDate() : date.getDate();
            var hours=date.getHours()<10 ? "0"+date.getHours() : date.getHours();
            var minutes=date.getMinutes()<10 ? "0"+date.getMinutes() : date.getMinutes();
            var seconds=date.getSeconds()<10 ? "0"+date.getSeconds() : date.getSeconds();
            // 拼接
            return year+"-"+month+"-"+day+" "+hours+":"+minutes+":"+seconds;
        }
    }
};
</script>

<style scoped>
</style>
