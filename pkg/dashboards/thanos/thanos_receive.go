package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"

	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withThanosReceiveRemoteWriteGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write v1 - Incoming Requests",
		panelgroup.PanelsPerLine(3),
		panels.RemoteWriteRequestRate(datasource, labelMatcher),
		panels.RemoteWriteRequestErrors(datasource, labelMatcher),
		panels.RemoteWriteRequestDurations(datasource, labelMatcher),
	)
}

func withThanosReceiveRemoteWriteTenantedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write v1 - Incoming Requests (tenanted)",
		panelgroup.PanelsPerLine(3),
		panels.TenantedRemoteWriteRequestRate(datasource, labelMatcher),
		panels.TenantedRemoteWriteRequestErrors(datasource, labelMatcher),
		panels.TenantedRemoteWriteRequestDurations(datasource, labelMatcher),
	)
}

func withThanosReceiveRemoteWriteHTTPGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write v1 - HTTP Requests",
		panelgroup.PanelsPerLine(3),
		panels.AvgRemoteWriteRequestSize(datasource, labelMatcher),
		panels.AvgFailedRemoteWriteRequestSize(datasource, labelMatcher),
		panels.InflightRemoteWriteRequests(datasource, labelMatcher),
	)
}

func withThanosReceiveRemoteWriteSeriesSampleGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Series and Samples (tenanted)",
		panelgroup.PanelsPerLine(4),
		panels.RemoteWriteSeriesRate(datasource, labelMatcher),
		panels.RemoteWriteSeriesNotWrittenRate(datasource, labelMatcher),
		panels.RemoteWriteSamplesRate(datasource, labelMatcher),
		panels.RemoteWriteSamplesNotWrittenRate(datasource, labelMatcher),
	)
}

func withThanosReceiveRemoteWriteReplicationGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write Replication",
		panelgroup.PanelsPerLine(2),
		panels.RemoteWriteReplicationRate(datasource, labelMatcher),
		panels.RemoteWriteReplicationErrorRate(datasource, labelMatcher),
	)
}

func withThanosReceiveRemoteWriteForwardGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write Forward",
		panelgroup.PanelsPerLine(2),
		panels.RemoteWriteForwardRate(datasource, labelMatcher),
		panels.RemoteWriteForwardErrorRate(datasource, labelMatcher),
	)
}

func withThanosReceiveWriteGRPCUnaryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Write gRPC Unary (WritableStore)",
		panelgroup.PanelsPerLine(3),
		panels.WriteGRPCUnaryRate(datasource, labelMatcher),
		panels.WriteGRPCUnaryErrors(datasource, labelMatcher),
		panels.WriteGPRCUnaryDurations(datasource, labelMatcher),
	)
}

func withPrometheusStorageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage",
		panelgroup.PanelsPerLine(3),
		panels.ReceiveAppendedSampleRate(datasource, labelMatcher),
		panels.ReceiveHeadSeries(datasource, labelMatcher),
		panels.ReceiveHeadChunks(datasource, labelMatcher),
	)
}

func BuildThanosReceiveOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("thanos-receive-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Thanos / Receive / Overview"),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("thanos_build_info{container=\"thanos-receive\"}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("job"),
				listVar.AllowMultiple(true),
			),
		),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "thanos_build_info"),
		dashboard.AddVariable("namespace",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("namespace",
					labelValuesVar.Matchers("thanos_status"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("namespace"),
			),
		),
		dashboard.AddVariable("tenant",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("tenant",
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"prometheus_tsdb_head_max_time{container=\"thanos-receive\"}",
							[]promql.LabelMatcher{clusterLabelMatcher, {Name: "job", Type: "=", Value: "$job"}, {Name: "namespace", Type: "=", Value: "$namespace"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("tenant"),
				listVar.AllowMultiple(true),
			),
		),
		withThanosReceiveRemoteWriteGroup(datasource, clusterLabelMatcher),
		withThanosReceiveRemoteWriteTenantedGroup(datasource, clusterLabelMatcher),
		withThanosReceiveRemoteWriteHTTPGroup(datasource, clusterLabelMatcher),
		withThanosReceiveRemoteWriteSeriesSampleGroup(datasource, clusterLabelMatcher),
		withThanosReceiveRemoteWriteReplicationGroup(datasource, clusterLabelMatcher),
		withThanosReceiveRemoteWriteForwardGroup(datasource, clusterLabelMatcher),
		withThanosReceiveWriteGRPCUnaryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcher),
		withThanosBucketUploadGroup(datasource, clusterLabelMatcher),
		withPrometheusStorageGroup(datasource, clusterLabelMatcher),
		withThanosResourcesGroup(datasource, clusterLabelMatcher),
	)
}
