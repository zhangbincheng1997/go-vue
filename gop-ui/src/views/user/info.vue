<template>
  <el-container>
    <el-form ref="form" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="角色" prop="role">
        <el-select v-model="form.role" placeholder="请选择角色">
          <el-option v-for="role in roleOptions" :key="role.id" :label="role.desc" :value="role.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="头像" prop="avatar">
        <el-upload
          class="avatar-uploader"
          action=""
          :http-request="uploadAvatar"
          :show-file-list="false"
          accept=".jpg, .jpeg, .png">
          <img v-if="form.avatar" :src="form.avatar" class="avatar">
          <i v-else class="el-icon-plus avatar-uploader-icon"></i>
        </el-upload>
      </el-form-item>
      <el-form-item label="昵称" prop="name">
        <el-input v-model="form.name" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="简介" prop="introduction">
        <el-input v-model="form.introduction" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm">提交</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>
  </el-container>
</template>

<script>
import { getInfo, updateInfo } from '@/api/user'
import { upload } from '@/api/upload'

export default {
  data() {
    return {
      form: {
        role: '',
        avatar: '',
        name: '',
        introduction: '',
      },
      rules: {
        role: [
          { required: true, message: '请输入角色', trigger: 'blur' }
        ]
      },
      roleOptions: [
        { id: 'admin', desc: '管理员' },
        { id: 'guest', desc: '游客' }
      ]
    }
  },
  created() {
    getInfo().then(res => {
      this.form = res.data
    })
  },
  methods: {
    uploadAvatar(param) {
      if (param.file.size > 1024 * 1024 * 10) {
        this.$message.error('上传图片大小不能超过10MB!')
        return
      }
      let formData = new FormData()
      formData.append('file', param.file)
      upload(formData).then(res => {
        this.$message({ type: 'success', message: res })
        this.form.avatar = res.data
      })
    },
    submitForm() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          updateInfo(this.form).then(res => {
            this.$message({ type: 'success', message: res })
          })
        } else {
          console.log('error submit!!')
          return false
        }
      });
    },
    resetForm() {
      this.$refs.form.resetFields();
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
