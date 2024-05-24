package controllers

import (
    "context"
    "time"

    batchv1 "k8s.io/api/batch/v1"
    batchv1beta1 "k8s.io/api/batch/v1beta1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    dariiav1 "github.com/kulik-dari/testing/api/dariia/v1"
)

// ClusterScanReconciler reconciles a ClusterScan object
type ClusterScanReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dariia.example.com,resources=clusterscans,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dariia.example.com,resources=clusterscans/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dariia.example.com,resources=clusterscans/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=jobs;cronjobs,verbs=get;list;watch;create;update;patch;delete

func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    // Fetch the ClusterScan instance
    var clusterScan dariiav1.ClusterScan
    if err := r.Get(ctx, req.NamespacedName, &clusterScan); err != nil {
        logger.Error(err, "unable to fetch ClusterScan")
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Determine if the ClusterScan should create a Job or a CronJob
    if clusterScan.Spec.Schedule == "" {
        // One-off execution: Create a Job
        job := &batchv1.Job{
            ObjectMeta: metav1.ObjectMeta{
                Name:      clusterScan.Name,
                Namespace: clusterScan.Namespace,
            },
            Spec: clusterScan.Spec.JobTemplate,
        }
        if err := ctrl.SetControllerReference(&clusterScan, job, r.Scheme); err != nil {
            return ctrl.Result{}, err
        }
        if err := r.Create(ctx, job); err != nil {
            return ctrl.Result{}, err
        }
        clusterScan.Status.JobName = job.Name
    } else {
        // Recurring execution: Create a CronJob
        cronJob := &batchv1beta1.CronJob{
            ObjectMeta: metav1.ObjectMeta{
                Name:      clusterScan.Name,
                Namespace: clusterScan.Namespace,
            },
            Spec: clusterScan.Spec.CronJobTemplate,
        }
        cronJob.Spec.Schedule = clusterScan.Spec.Schedule
        cronJob.Spec.JobTemplate.Spec = clusterScan.Spec.JobTemplate
        if err := ctrl.SetControllerReference(&clusterScan, cronJob, r.Scheme); err != nil {
            return ctrl.Result{}, err
        }
        if err := r.Create(ctx, cronJob); err != nil {
            return ctrl.Result{}, err
        }
        clusterScan.Status.JobName = cronJob.Name
    }

    // Update the ClusterScan status with the job information
    clusterScan.Status.LastExecutionTime = &metav1.Time{Time: time.Now()}
    // TODO: Update LastExecutionResult based on the job's execution result

    if err := r.Status().Update(ctx, &clusterScan); err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}

func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&dariiav1.ClusterScan{}).
        Complete(r)
}

