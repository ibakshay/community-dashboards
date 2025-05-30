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

func withThanosBlockOperationsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Block Operations",
		panelgroup.PanelsPerLine(4),
		panels.BlockLoadRate(datasource, labelMatcher),
		panels.BlockLoadErrors(datasource, labelMatcher),
		panels.BlockDropRate(datasource, labelMatcher),
		panels.BlockDropErrors(datasource, labelMatcher),
	)
}

func withThanosCacheOperationsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Cache Operations",
		panelgroup.PanelsPerLine(4),
		panels.CacheRequestRate(datasource, labelMatcher),
		panels.CacheHitRate(datasource, labelMatcher),
		panels.CacheItemsAddRate(datasource, labelMatcher),
		panels.CacheItemsEvictRate(datasource, labelMatcher),
	)
}

func withThanosQueryOperationsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Operations",
		panelgroup.PanelsPerLine(4),
		panels.BlocksQueried(datasource, labelMatcher),
		panels.DataFetched(datasource, labelMatcher),
		panels.DataTouched(datasource, labelMatcher),
		panels.ResultSeries(datasource, labelMatcher),
	)
}

func withThanosQueryOperationDurationGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Operation Durations",
		panelgroup.PanelsPerLine(3),
		panels.GetAllSeriesDurations(datasource, labelMatcher),
		panels.MergeDurations(datasource, labelMatcher),
		panels.GateWaitingDurations(datasource, labelMatcher),
	)
}

func withThanosStoreSentGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Store Sent Chunk Size",
		panelgroup.PanelsPerLine(1),
		panels.StoreSentChunkSizes(datasource, labelMatcher),
	)
}

func BuildThanosStoreOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("thanos-store-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Thanos / Store Gateway / Overview"),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("thanos_build_info{container=\"thanos-store\"}"),
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
		withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcher),
		withThanosBucketOperationsGroup(datasource, clusterLabelMatcher),
		withThanosBlockOperationsGroup(datasource, clusterLabelMatcher),
		withThanosCacheOperationsGroup(datasource, clusterLabelMatcher),
		withThanosQueryOperationsGroup(datasource, clusterLabelMatcher),
		withThanosQueryOperationDurationGroup(datasource, clusterLabelMatcher),
		withThanosStoreSentGroup(datasource, clusterLabelMatcher),
		withThanosResourcesGroup(datasource, clusterLabelMatcher),
	)
}
