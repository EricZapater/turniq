import { createRouter, createWebHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";
import { useAuthStore } from "../stores/auth.store";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/login",
      name: "login",
      component: LoginView,
      meta: { public: true },
    },
    {
      path: "/",
      component: () => import("../layouts/MainLayout.vue"),
      children: [
        {
          path: "",
          name: "home",
          component: HomeView,
        },
        {
          path: "customers",
          name: "customers-list",
          component: () => import("../views/customers/CustomersListView.vue"),
        },
        {
          path: "customers/new",
          name: "customers-create",
          component: () => import("../views/customers/CustomerDetailView.vue"),
        },
        {
          path: "customers/:id",
          name: "customers-view",
          component: () => import("../views/customers/CustomerDetailView.vue"),
        },
        {
          path: "customers/:id/edit",
          name: "customers-edit",
          component: () => import("../views/customers/CustomerDetailView.vue"),
        },
      ],
    },
    {
      path: "/",
      component: () => import("../layouts/MainLayout.vue"),
      children: [
        {
          path: "my-customer",
          name: "my-customer",
          component: () => import("../views/MyProfileView.vue"),
        },
        {
          path: "company-users",
          name: "company-users",
          component: () => import("../views/CompanyUsersView.vue"),
        },
        // Shopfloors
        {
          path: "shopfloors",
          name: "shopfloors-list",
          component: () => import("../views/shopfloors/ShopfloorsListView.vue"),
        },
        {
          path: "shopfloors/new",
          name: "shopfloors-create",
          component: () =>
            import("../views/shopfloors/ShopfloorDetailView.vue"),
        },
        {
          path: "shopfloors/:id",
          name: "shopfloors-detail",
          component: () =>
            import("../views/shopfloors/ShopfloorDetailView.vue"),
        },
        {
          path: "shopfloors/:id/edit",
          name: "shopfloors-edit",
          component: () =>
            import("../views/shopfloors/ShopfloorDetailView.vue"),
        },
        // Workcenters
        {
          path: "workcenters",
          name: "workcenters-list",
          component: () =>
            import("../views/workcenters/WorkcentersListView.vue"),
        },
        {
          path: "workcenters/new",
          name: "workcenters-create",
          component: () =>
            import("../views/workcenters/WorkcenterDetailView.vue"),
        },
        {
          path: "workcenters/:id",
          name: "workcenters-detail",
          component: () =>
            import("../views/workcenters/WorkcenterDetailView.vue"),
        },
        {
          path: "workcenters/:id/edit",
          name: "workcenters-edit",
          component: () =>
            import("../views/workcenters/WorkcenterDetailView.vue"),
        },
        // Operators
        {
          path: "operators",
          name: "operators-list",
          component: () => import("../views/operators/OperatorsListView.vue"),
        },
        {
          path: "operators/new",
          name: "operators-create",
          component: () => import("../views/operators/OperatorDetailView.vue"),
        },
        {
          path: "operators/:id",
          name: "operators-detail",
          component: () => import("../views/operators/OperatorDetailView.vue"),
        },
        {
          path: "operators/:id/edit",
          name: "operators-edit",
          component: () => import("../views/operators/OperatorDetailView.vue"),
        },
        // Jobs
        {
          path: "jobs",
          name: "jobs-list",
          component: () => import("../views/jobs/JobsListView.vue"),
        },
        {
          path: "jobs/new",
          name: "jobs-create",
          component: () => import("../views/jobs/JobDetailView.vue"),
        },
        {
          path: "jobs/:id",
          name: "jobs-detail",
          component: () => import("../views/jobs/JobDetailView.vue"),
        },
        {
          path: "jobs/:id/edit",
          name: "jobs-edit",
          component: () => import("../views/jobs/JobDetailView.vue"),
          meta: { requiresAuth: true },
        },
        // Planning
        {
          path: "planning",
          name: "planning",
          component: () =>
            import("../views/planning/ShopfloorPlanningView.vue"),
          meta: { requiresAuth: true },
        },
        // Shifts
        {
          path: "shifts",
          name: "shifts-list",
          component: () => import("../views/shifts/ShiftsListView.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "shifts/new",
          name: "shifts-create",
          component: () => import("../views/shifts/ShiftDetailView.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "shifts/:id",
          name: "shifts-detail",
          component: () => import("../views/shifts/ShiftDetailView.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "shifts/:id/edit",
          name: "shifts-edit",
          component: () => import("../views/shifts/ShiftDetailView.vue"),
          meta: { requiresAuth: true },
        },
        // Reports
        {
          path: "reports/billing",
          name: "reports-billing",
          component: () => import("../views/reports/BillingReportView.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "reports/schedule",
          name: "reports-schedule",
          component: () => import("../views/reports/ScheduleReportView.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "reports/time-entries",
          name: "reports-time-entries",
          component: () => import("../views/reports/TimeEntryReportView.vue"),
          meta: { requiresAuth: true },
        },
      ],
    },
  ],
});

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();

  if (!to.meta.public && !authStore.isAuthenticated) {
    next("/login");
  } else if (to.path === "/login" && authStore.isAuthenticated) {
    next("/");
  } else {
    next();
  }
});

export default router;
