<template>
  <el-container>
    <el-form ref="form" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="头像" prop="avatar">
        <el-upload
          class="avatar-uploader"
          action=""
          :http-request="uploadAvatar"
          :show-file-list="false"
          accept=".jpg, .jpeg, .png"
        >
          <img v-if="form.avatar" :src="form.avatar" class="avatar">
          <i v-else class="el-icon-plus avatar-uploader-icon" />
        </el-upload>
      </el-form-item>
      <el-form-item label="昵称" prop="name">
        <el-input v-model="form.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="简介" prop="introduction">
        <el-input v-model="form.introduction" autocomplete="off" />
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
        avatar: '',
        name: '',
        introduction: ''
      },
      rules: {
        avatar: [
          { required: true, message: '请上传头像', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入昵称', trigger: 'blur' }
        ],
        introduction: [
          { required: true, message: '请输入简介', trigger: 'blur' }
        ]
      }
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
      const formData = new FormData()
      formData.append('file', param.file)
      upload(formData).then(res => {
        this.form.avatar = res.data
        this.$message({ type: 'success', message: res })
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
      })
    },
    resetForm() {
      this.$refs.form.resetFields()
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
