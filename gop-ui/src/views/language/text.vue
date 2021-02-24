<template>
  <div class="app-container">
    <div class="filter-container">
      <el-select v-model="listQuery.status" placeholder="状态" clearable style="width: 150px">
        <el-option v-for="status in statusOptions" :key="status.id" :label="status.desc" :value="status.id" />
      </el-select>
      <el-input v-model="listQuery.keyword" placeholder="请输入ID/中文" prefix-icon="el-icon-search" clearable style="width: 200px" />
      <el-button type="primary" icon="el-icon-search" @click="getList">
        查询
      </el-button>
      <!-- <el-button style="float:right;" type="primary" icon="el-icon-edit" @click="handleCreate">
        添加
      </el-button> -->
      <el-button :loading="exportLoading" style="float:right;" type="primary" icon="el-icon-download" @click="handleExport">
        导出
      </el-button>
      <el-button :loading="importLoading" style="float:right;" type="primary" icon="el-icon-upload" @click="importVisible = true">
        导入
      </el-button>
      <el-button v-if="selectIds && selectIds.length > 0" style="float:right;" type="success" icon="el-icon-check" @click="selectVisible = true">
        设置状态
      </el-button>
      <el-button v-if="selectIds && selectIds.length > 0" style="float:right;" type="danger" icon="el-icon-delete" @click="handleRemove()">
        删除
      </el-button>
    </div>
    <el-table
      v-loading="listLoading"
      :data="list"
      border
      highlight-current-row
      style="width: 100%;"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
      @cell-click="handleCellClick"
    >
      <el-table-column type="selection" width="40" />
      <el-table-column label="ID" align="center" prop="id" width="100" sortable="custom" fixed="left" />
      <el-table-column label="文本" align="center" prop="text">
        <template scope="scope">
          <span v-if="scope.row.id === editId">
            <el-input size="small" v-model="scope.row.text" @blur="handleInputBlur(scope)" type="textarea" autosize maxlength="1000" show-word-limit />
          </span>
          <span v-else>{{ scope.row.text }}</span>
        </template>
      </el-table-column>
      <el-table-column label="翻译" align="center" prop="text2">
        <template scope="scope">
          <span v-if="scope.row.id === editId2">
            <el-input size="small" v-model="scope.row.text2" @blur="handleInputBlur(scope)" type="textarea" autosize maxlength="1000" show-word-limit />
          </span>
          <span v-else>{{ scope.row.text2 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="120" fixed="right">
        <template slot-scope="scope">
          <el-button type="text" size="small" @click="handleEdit(scope.row)">翻译</el-button>
          <el-button type="text" size="small" @click="handleHistory(scope.row.id)">历史</el-button>
        </template>
      </el-table-column>
      <el-table-column label="状态" align="center" width="120" fixed="right">
        <template slot-scope="scope">
          <el-select v-model="scope.row.status" placeholder="请选择" size="mini" @change="handleStatus(scope.row)">
            <el-option v-for="status in statusOptions" :key="status.id" :label="status.desc" :value="status.id" />
          </el-select>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog
      title="翻译"
      :visible.sync="editVisible"
      :close-on-click-modal="false"
      @close="resetEditForm"
    >
      <el-form ref="item" :model="item" label-width="120px">
        <el-form-item label="路径" prop="path">
          <el-input v-model="item.path" disabled />
        </el-form-item>
        <el-form-item label="变量名" prop="property">
          <el-input v-model="item.property" disabled />
        </el-form-item>
        <el-form-item label="行数" prop="line">
          <el-input v-model="item.line" disabled />
        </el-form-item>
        <el-form-item label="文本" prop="text">
          <el-input v-model="item.text" type="textarea" autosize disabled />
        </el-form-item>
        <el-form-item label="翻译" prop="text2" :rules="{ required: true, message: '请输入中文', trigger: 'blur' }">
          <el-input v-model="item.text2" type="textarea" autosize maxlength="1000" show-word-limit />
        </el-form-item>
        <el-form-item label="备注" prop="comment">
          <el-input v-model="item.comment" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitEditForm">确 定</el-button>
        <el-button @click="resetEditForm">取 消</el-button>
      </div>
    </el-dialog>

    <el-dialog
      title="导入"
      :visible.sync="importVisible"
      width="30%"
      center
      @close="resetImportForm"
    >
      <el-form>
        <el-form-item label="" align="center" prop="file">
          <el-upload
            ref="upload"
            drag
            action=""
            accept=".csv"
            :auto-upload="false"
            :limit="1">
            <i class="el-icon-upload"></i>
            <div>将.csv文件拖到此处，或<em>点击上传</em></div>
          </el-upload>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitImportForm">确 定</el-button>
        <el-button @click="resetImportForm">取 消</el-button>
      </div>
    </el-dialog>

    <el-dialog
      title="设置状态"
      :visible.sync="selectVisible"
      width="30%"
      center
      @close="resetSelectForm"
    >
      <el-form>
        <el-form-item align="center">
          <el-select v-model="selectStatus" placeholder="请选择">
            <el-option v-for="status in statusOptions" :key="status.id" :label="status.desc" :value="status.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitSelectForm">确 定</el-button>
        <el-button @click="resetSelectForm">取 消</el-button>
      </div>
    </el-dialog>

    <el-dialog
      title="历史"
      :visible.sync="historyVisible"
      @close="historyVisible = false"
    >
      <el-table
        :data="historyData.slice((historyPage-1)*historyLimit,historyPage*historyLimit)"
        border
        highlight-current-row
        style="width: 100%"
      >
        <el-table-column label="翻译" align="center" prop="text" />
        <el-table-column label="备注" align="center" prop="comment" width="200" />
        <el-table-column label="操作者" align="center" prop="operator" width="100" />
        <el-table-column label="时间" align="center" width="100">
          <template slot-scope="scope">
            <span> {{ scope.row.time | parseTime }} </span>
          </template>
        </el-table-column>
      </el-table>
      <pagination v-show="historyData.length>0" :total="historyData.length" :page.sync="historyPage" :limit.sync="historyLimit" />
    </el-dialog>
  </div>
</template>

<script>
import { getStatusOptions, getList, add, remove, update, updateStatus, updateText, updateRecordText, getRecordList, importData, exportText } from '@/api/language'
import { parseTime } from '@/utils'

import axios from 'axios'
import Pagination from '@/components/Pagination'

// 编辑
const defaultItem = {
  id: undefined, // ID
  path: undefined, // 路径
  property: undefined, // 变量名
  line: undefined, // 行数
  text: undefined, // 文本
  text2: undefined, // 翻译
  comment: undefined // 备注
}

// 查询
const defaultListQuery = {
  page: 1,
  limit: 10,
  language: undefined, // 语言：korea, tradition
  status: undefined, // 状态
  sort: undefined, // ID排序
  keyword: undefined // ID/中文
}

export default {
  components: { Pagination },
  props: ['table', 'language'], // router传参
  data() {
    return {
      listLoading: false,
      listQuery: Object.assign({}, defaultListQuery),
      list: [],
      total: 0,
      item: Object.assign({}, defaultItem),

      selectIds: undefined,
      selectStatus: undefined,
      selectVisible: false, // 多选修改状态

      editVisible: false, // 编辑
      historyVisible: false, // 历史
      importVisible: false, // 导入
      importLoading: false,
      exportLoading: false,

      editId: -1, // 记录正在编辑text的rowId
      editId2: -1, // 记录正在编辑text2的rowId
      editContent: undefined, // 编辑前的内容

      historyData: [],
      historyPage: 1,
      historyLimit: 5,

      statusOptions: undefined,
      statusMap: {}
    }
  },
  created() {
    this.listQuery.table = this.table
    this.listQuery.language = this.language
    this.init()
  },
  methods: {
    init() {
      getStatusOptions().then(res => {
        this.statusOptions = res.data
        for(const status of this.statusOptions) this.statusMap[status.id] = status.desc
      })
      this.getList()
    },
    getList() {
      this.listLoading = true
      getList(this.listQuery).then(res => {
        this.total = res.data.count
        this.list = []
        if (res.data.list instanceof Array) {
          var lst = res.data.list
          // 过滤语言类型数据
          for (var i in lst) {
            var item = lst[i]
            item.text2 = typeof(item[this.language]) === 'undefined' ? '' :  item[this.language].text
            item.status = typeof(item[this.language]) === 'undefined' ? -1 :  item[this.language].status // NONE 特殊处理
            this.list.push(item)
          }
         }
      }).finally(() => {
        this.listLoading = false
      })
    },
    handleEdit(row) {
      this.editVisible = true
      this.item = JSON.parse(JSON.stringify(row)) // 深拷贝
    },
    submitEditForm() {
      this.$refs.item.validate((valid) => {
        if (!valid) return false

        var _item = (({id, text2, comment}) => ({id, text: text2, comment}))(this.item) // pick 2->1
        updateRecord(this.table, this.language, _item).then(res => {
          this.$message({ type: 'success', message: '操作成功：' + JSON.stringify(_item) })
          this.resetEditForm()
          this.getList()
        })
      })
    },
    resetEditForm() {
      this.editVisible = false
      this.$refs.item.resetFields();
      this.item = Object.assign({}, defaultItem)
    },
    handleRemove() {
      this.$confirm('是否删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(this.table, this.language, this.selectIds).then(res => {
          this.$message({ type: 'success', message: '删除成功：' + this.selectIds })
          this.getList()
        })
      })
    },
    handleCellClick(row, column, cell, event) {
      this.editId = -1
      this.editId2 = -1
      if(column.property === 'text') {
        this.editId = row.id
        this.editContent = row.text
      }
      if(column.property === 'text2') {
        this.editId2 = row.id
        this.editContent = row.text2
      }
    },
    handleInputBlur(scope) {
      if(this.editId !== -1) {
        let text = scope.row.text
        if(scope.row.text === this.editContent) return; // 没有修改
        updateText(this.table, this.language, scope.row.id, text).then(res => {
          this.$message({ type: 'success', message: '修改成功：' + JSON.stringify(text) })
        })
      }
      if(this.editId2 !== -1) {
        let text = scope.row.text2
        if(text === this.editContent) return; // 没有修改
        updateRecordText(this.table, this.language, scope.row.id, text).then(res => {
          this.$message({ type: 'success', message: '修改成功：' + JSON.stringify(text) })
        })
      }
      this.editId = -1
      this.editId2 = -1
      this.editContent = undefined
    },
    handleSortChange({ column, prop, order }) {
      this.listQuery.sort = order === 'ascending' ? 1 : -1
      this.getList()
    },
    handleSelectionChange(val) {
      this.selectIds = []
      for (let i = 0; i < val.length; i++) if (val[i].id) this.selectIds.push(val[i].id)
    },
    submitSelectForm() {
      if (typeof(this.selectIds) === 'undefined' || typeof(this.selectStatus) === 'undefined') return
      updateStatus(this.table, this.language, this.selectIds, this.selectStatus).then(res => {
        this.$message({ type: 'success', message: '修改状态成功：' + this.statusMap[this.selectStatus] })
        this.resetSelectForm()
        this.getList()
      })
    },
    resetSelectForm() {
      this.selectVisible = false
      this.selectStatus = undefined
    },
    handleStatus(row) {
      updateStatus(this.table, this.language, [row.id], row.status).then(res => {
        this.$message({ type: 'success', message: '修改状态成功：' + this.statusMap[row.status] })
      })
    },
    handleHistory(id) {
      this.historyVisible = true
      getRecordList(this.table, this.language, id).then(res => {
        if (res.data instanceof Array) {
          this.historyData = res.data
        } else {
          this.historyData = []
        }
      })
    },
    submitImportForm() {
      if(this.$refs.upload.uploadFiles.length === 0) return
      let formData = new FormData()
      formData.append('file', this.$refs.upload.uploadFiles[0].raw)
      formData.append('table',this.table)
      formData.append('language', this.language)
      this.resetImportForm()
      this.importLoading = true
      importData(formData).then(res => {
        this.$message({ type: 'success', message: res.data })
      }).finally(() => {
        this.importLoading = false
      })
    },
    resetImportForm() {
      this.importVisible = false
      this.$refs.upload.clearFiles()
    },
    handleExport() {
      this.$confirm('是否导出所有【已完成】？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.exportLoading = true
        axios({
          method: 'get',
          baseURL: process.env.VUE_APP_BASE_API,
          url: '/item/export',
          params: {
            table: this.table,
            language: this.language,
          },
          responseType: 'blob'
        }).then(response => {
          this.download(response)
        }).finally(() => {
          this.exportLoading = false
        })
      })
    },
    // 下载文件
    download (res) {
      let link = document.createElement('a')
      link.href = window.URL.createObjectURL(new Blob([res.data]))
      link.download = res.headers['content-disposition'].split('=')[1]
      link.style.display = 'none'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }
  },
  filters: {
    parseTime(time) {
      return parseTime(time, '{y}-{m}-{d} {h}:{i}:{s}')
    }
  }
}
</script>
