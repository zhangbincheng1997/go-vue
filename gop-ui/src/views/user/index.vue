<template>
  <div class="app-container">
    <div class="filter-container">
      <el-select v-model="listQuery.role" placeholder="角色" clearable style="width: 150px">
        <el-option v-for="role in roleOptions" :key="role.id" :label="role.desc" :value="role.id" />
      </el-select>
      <el-button type="primary" icon="el-icon-search" @click="getList">
        查询
      </el-button>
      <el-button style="float:right;" type="primary" icon="el-icon-download" @click="visible = true">
        注册
      </el-button>
    </div>
    <el-table
      v-loading="listLoading"
      :data="list"
      border
      highlight-current-row
      style="width: 100%"
      @sort-change="handleSortChange"
    >
      <el-table-column label="ID" align="center" prop="id" width="100" sortable="custom" fixed="left" />
      <el-table-column label="用户名" align="center" prop="username" width="200" />
      <el-table-column label="头像" align="center" prop="avatar" width="200">
        <template slot-scope="scope">
          <el-image style="width: 100px; height: 100px" :src="scope.row.avatar"></el-image>
        </template>
      </el-table-column>
      <el-table-column label="昵称" align="center" prop="name" width="200" />
      <el-table-column label="简介" align="center" prop="introduction" />
      <el-table-column label="角色" align="center" prop="role" width="200">
        <template slot-scope="scope">
          <el-select v-model="scope.row.role" placeholder="请选择" size="mini" @change="handleRole(scope.row)">
            <el-option v-for="role in roleOptions" :key="role.id" :label="role.desc" :value="role.id" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="120" fixed="right">
        <template slot-scope="scope">
          <el-button @click="handleDelete(scope.row)" type="text" size="small">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  
    <el-dialog
      title="注册"
      :visible.sync="visible"
      width="30%"
      center
      @close="resetForm"
    >
      <el-form ref="form" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="用户名" props="username">
          <el-input v-model="form.username"></el-input>
        </el-form-item>
        <el-form-item label="密码" props="password">
          <el-input v-model="form.password"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="resetForm">取 消</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { register, getUserList, updateRole, deleteUser } from '@/api/user'

import Pagination from '@/components/Pagination'

// 查询
const defaultListQuery = {
  page: 1,
  limit: 10,
  sort: undefined, // ID排序
  role: undefined // 角色
}

export default {
  components: { Pagination },
  data() {
    return {
      listLoading: false,
      listQuery: Object.assign({}, defaultListQuery),
      list: [],
      total: 0,

      visible: false,
      form: {
        username: '',
        password: '',
      },
      rules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ]
      },

      roleOptions: [
        { id: 'admin', desc: '管理员' },
        { id: 'guest', desc: '游客' }
      ],
      roleMap: {
        'admin': '管理员',
        'guest': '游客'
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      getUserList(this.listQuery).then(res => {
        this.total = res.data.count
        this.list = []
        if (res.data.list instanceof Array) {
          this.list = res.data.list
        }
      }).finally(() => {
        this.listLoading = false
      })
    },
    handleSortChange({ column, prop, order }) {
      this.listQuery.sort = order === 'ascending' ? true : false
      this.getList()
    },
    submitForm() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          register(this.form).then(res => {
            this.$message({ type: 'success', message: res })
            this.resetForm()
          })
        } else {
          console.log('error submit!!')
          return false
        }
      });
    },
    resetForm() {
      this.$refs.form.resetFields()
      this.visible = false
    },
    handleDelete(row) {
      this.$confirm('是否删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteUser(row.id).then(res => {
          this.$message({ type: 'success', message: res })
          this.getList()
        })
      })
    },
    handleRole(row) {
      updateRole(row.id, row.role).then(res => {
          this.$message({ type: 'success', message: res })
      })
    }
  }
}
</script>

<style lang="scss">
.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style>
