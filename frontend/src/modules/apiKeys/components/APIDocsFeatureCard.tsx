/**
 * @file APIDocsFeatureCard.tsx
 * @brief Karta skupiny API funkcí používaná v přehledu dokumentace.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";

import type { ApiFeatureDoc } from "../data/apiDocs";

type Props = {
  feature: ApiFeatureDoc;
  selected: boolean;
  onClick: () => void;
};

export default function ApiDocsFeatureCard({
  feature,
  selected,
  onClick,
}: Props) {
  return (
    <Card
      sx={{
        flexShrink: 0,
        height: "100%",
        overflow: "hidden",
        backgroundColor: selected ? "var(--border-input)" : "var(--bg-card)",
        "&:hover": {
          backgroundColor: selected ? "var(--border-input)" : "var(--border)",
        },
      }}
    >
      <CardActionArea
        onClick={onClick}
        sx={{
          height: "100%",
        }}
      >
        <CardContent
          sx={{
            py: 1.5,
            height: "100%",
            overflow: "hidden",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
          }}
        >
          <Typography variant="body1" className="form-label" sx={{ mb: 0.5 }}>
            {feature.title}
          </Typography>

          <Typography
            variant="caption"
            className="form-label"
            sx={{
              opacity: 0.8,
              overflow: "hidden",
              textOverflow: "ellipsis",
              display: "-webkit-box",
              WebkitBoxOrient: "vertical",
              WebkitLineClamp: 2,
            }}
          >
            {feature.summary}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
