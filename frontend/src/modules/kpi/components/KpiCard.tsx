import { useNavigate } from "react-router-dom";
import type { KpiDefinitionsQuery } from "../../../generated/graphql";

import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Chip from "@mui/material/Chip";

type Props = {
  kpi: KpiDefinitionsQuery["kpiDefinitions"][number];
  sdTypeMap: Map<string, string>;
};

export default function KpiCard({ kpi, sdTypeMap }: Props) {
  const navigate = useNavigate();

  const typeLabel =
    kpi.sdTypeUID || sdTypeMap.get(kpi.sdTypeID) || kpi.sdTypeID || "\u00A0";

  return (
    <Card
      sx={{
        height: "100%",
        minHeight: 120,
        backgroundColor: "var(--bg-card)",
      }}
    >
      <CardActionArea
        onClick={() => navigate(`/kpi/${kpi.id}`)}
        sx={{
          height: "100%",
          display: "flex",
          alignItems: "stretch",
        }}
      >
        <CardContent
          sx={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "space-between",
            width: "100%",
          }}
        >
          {/* TITLE */}
          <Typography variant="h6" className="color-label">
            {kpi.label}
          </Typography>

          {/* META */}
          <Box sx={{ mt: 1 }}>
            <Typography
              variant="body2"
              sx={{ opacity: 0.7 }}
              className="color-label"
            >
              {typeLabel}
            </Typography>

            <Box sx={{ mt: 1 }}>
              <Chip label={kpi.sdInstanceMode} size="small" color="primary" />
            </Box>
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
