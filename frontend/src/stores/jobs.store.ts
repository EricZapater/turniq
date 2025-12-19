import { defineStore } from "pinia";
import { ref } from "vue";
import {
  jobsApi,
  type Job,
  type JobRequest,
  type JobListParams,
} from "../api/jobs.api";

export const useJobsStore = defineStore("jobs", () => {
  const jobs = ref<Job[]>([]);
  const total = ref<number>(0);
  const currentJob = ref<Job | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);

  async function fetchJobs(params: JobListParams) {
    loading.value = true;
    error.value = null;
    try {
      const response = await jobsApi.list(params);
      jobs.value = response.data;
      total.value = response.data ? response.data.length : 0;
    } catch (err: any) {
      console.error("Error fetching jobs", err);
      error.value = err.message || "Error carregant feines";
    } finally {
      loading.value = false;
    }
  }

  async function fetchJob(id: string) {
    loading.value = true;
    error.value = null;
    currentJob.value = null;
    try {
      const response = await jobsApi.get(id);
      currentJob.value = response.data;
    } catch (err: any) {
      console.error("Error fetching job", err);
      error.value = err.message || "Error carregant feina";
    } finally {
      loading.value = false;
    }
  }

  async function createJob(data: JobRequest) {
    loading.value = true;
    error.value = null;
    try {
      await jobsApi.create(data);
    } catch (err: any) {
      console.error("Error creating job", err);
      const msg =
        err.response?.data?.error || err.message || "Error creant feina";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  async function updateJob(id: string, data: JobRequest) {
    loading.value = true;
    error.value = null;
    try {
      await jobsApi.update(id, data);
    } catch (err: any) {
      console.error("Error updating job", err);
      error.value = err.message || "Error actualitzant feina";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function deleteJob(id: string) {
    loading.value = true;
    error.value = null;
    try {
      await jobsApi.delete(id);
    } catch (err: any) {
      console.error("Error deleting job", err);
      error.value = err.message || "Error eliminant feina";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    jobs,
    total,
    currentJob,
    loading,
    error,
    fetchJobs,
    fetchJob,
    createJob,
    updateJob,
    deleteJob,
  };
});
