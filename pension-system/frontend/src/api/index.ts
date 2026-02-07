import * as AuthController from '../../wailsjs/go/controllers/AuthController'
import * as DataController from '../../wailsjs/go/controllers/DataController'
import * as SurveyController from '../../wailsjs/go/controllers/SurveyController'
import * as IssueController from '../../wailsjs/go/controllers/IssueController'
import * as AccountController from '../../wailsjs/go/controllers/AccountController'

class ApiClient {
  // Auth APIs
  async login(username: string, password: string) {
    const result = await AuthController.Login(username, password)
    return JSON.parse(result)
  }

  async register(username: string, password: string, name: string) {
    const result = await AuthController.Register(username, password, name)
    return JSON.parse(result)
  }

  async getCurrentUser(userId: number) {
    const result = await AuthController.GetCurrentUser(userId)
    return JSON.parse(result)
  }

  // Data APIs
  async getElderlyList(page: number, pageSize: number, keyword: string = '') {
    const result = await DataController.GetElderlyList(page, pageSize, keyword)
    return JSON.parse(result)
  }

  async getElderly(id: number) {
    const result = await DataController.GetElderly(id)
    return JSON.parse(result)
  }

  async createElderly(data: any) {
    const result = await DataController.CreateElderly(JSON.stringify(data))
    return JSON.parse(result)
  }

  async updateElderly(id: number, data: any) {
    const result = await DataController.UpdateElderly(id, JSON.stringify(data))
    return JSON.parse(result)
  }

  async deleteElderly(id: number) {
    const result = await DataController.DeleteElderly(id)
    return JSON.parse(result)
  }

  async getSummary() {
    const result = await DataController.GetSummary()
    return JSON.parse(result)
  }

  // Survey APIs
  async getSurveyList(page: number, pageSize: number) {
    const result = await SurveyController.GetSurveyList(page, pageSize)
    return JSON.parse(result)
  }

  async getActiveSurveys() {
    const result = await SurveyController.GetActiveSurveys()
    return JSON.parse(result)
  }

  async getSurvey(id: number) {
    const result = await SurveyController.GetSurvey(id)
    return JSON.parse(result)
  }

  async createSurvey(title: string, description: string, options: string, startDate: string, endDate: string, createdBy: number) {
    const result = await SurveyController.CreateSurvey(title, description, options, startDate, endDate, createdBy)
    return JSON.parse(result)
  }

  async updateSurveyStatus(id: number, status: string) {
    const result = await SurveyController.UpdateSurveyStatus(id, status)
    return JSON.parse(result)
  }

  async vote(surveyId: number, option: string, userId: number) {
    const result = await SurveyController.Vote(surveyId, option, userId)
    return JSON.parse(result)
  }

  async getSurveyResults(surveyId: number) {
    const result = await SurveyController.GetSurveyResults(surveyId)
    return JSON.parse(result)
  }

  async deleteSurvey(id: number) {
    const result = await SurveyController.DeleteSurvey(id)
    return JSON.parse(result)
  }

  // Issue APIs
  async getIssueList(page: number, pageSize: number, status: string = 'all') {
    const result = await IssueController.GetIssueList(page, pageSize, status)
    return JSON.parse(result)
  }

  async getIssue(id: number) {
    const result = await IssueController.GetIssue(id)
    return JSON.parse(result)
  }

  async createIssue(title: string, description: string, category: string, priority: string, submitterId: number) {
    const result = await IssueController.CreateIssue(title, description, category, priority, submitterId)
    return JSON.parse(result)
  }

  async updateIssueStatus(id: number, status: string) {
    const result = await IssueController.UpdateIssueStatus(id, status)
    return JSON.parse(result)
  }

  async respondToIssue(id: number, response: string, handlerId: number) {
    const result = await IssueController.RespondToIssue(id, response, handlerId)
    return JSON.parse(result)
  }

  async closeIssue(id: number) {
    const result = await IssueController.CloseIssue(id)
    return JSON.parse(result)
  }

  async deleteIssue(id: number) {
    const result = await IssueController.DeleteIssue(id)
    return JSON.parse(result)
  }

  // Account APIs
  async updateName(newName: string, password: string, userId: number) {
    const result = await AccountController.UpdateName(newName, password, userId)
    return JSON.parse(result)
  }

  async updatePassword(currentPassword: string, newPassword: string, userId: number) {
    const result = await AccountController.UpdatePassword(currentPassword, newPassword, userId)
    return JSON.parse(result)
  }

  async deleteAccount(password: string, userId: number) {
    const result = await AccountController.DeleteAccount(password, userId)
    return JSON.parse(result)
  }

  async getUserList(page: number, pageSize: number) {
    const result = await AccountController.GetUserList(page, pageSize)
    return JSON.parse(result)
  }

  async deleteUserAccount(targetUserId: number, adminId: number) {
    const result = await AccountController.DeleteUserAccount(targetUserId, adminId)
    return JSON.parse(result)
  }

  async resetUserPassword(userId: number) {
    const result = await AccountController.ResetUserPassword(userId)
    return JSON.parse(result)
  }
}

export const api = new ApiClient()
export default api
