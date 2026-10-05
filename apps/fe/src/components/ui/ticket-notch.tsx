"use client";

import React from "react";
import { cn } from "@/lib/utils";

interface TicketNotchDividerProps {
  className?: string;
  fillColor?: string;
  borderColor?: string;
  dashedLineColor?: string;
  notchRadius?: number;
}

/**
 * TicketNotchDivider
 * 
 * Renders a true concave cutout notch on the left and right edges with a perforated dashed line in between.
 * Outside the notch arc is 100% transparent vector space (no fake solid background discs).
 * The 1px border stroke follows the inward curve seamlessly connecting to top and bottom card borders.
 */
export function TicketNotchDivider({
  className,
  fillColor = "#fdfcf9",
  borderColor = "#ded4c1",
  dashedLineColor = "#c8beaf",
  notchRadius = 10,
}: TicketNotchDividerProps) {
  const height = notchRadius * 2;
  const width = notchRadius + 4; // 14px for radius 10

  return (
    <div
      className={cn("relative flex items-center justify-between w-full select-none", className)}
      style={{ height: `${height}px` }}
      aria-hidden="true"
    >
      {/* Left Inward Notch SVG */}
      <svg
        width={width}
        height={height}
        viewBox={`0 0 ${width} ${height}`}
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="shrink-0 -ml-[1px] block"
      >
        {/* Paper Fill inside card (everything right of the concave arc) */}
        <path
          d={`M 0 0 L ${width} 0 L ${width} ${height} L 0 ${height} A ${notchRadius} ${notchRadius} 0 0 0 0 0 Z`}
          fill={fillColor}
        />
        {/* Inward border curve connecting top and bottom card edges */}
        <path
          d={`M 0.5 0 A ${notchRadius - 0.5} ${notchRadius - 0.5} 0 0 1 0.5 ${height}`}
          stroke={borderColor}
          strokeWidth="1"
          fill="none"
        />
      </svg>

      {/* Perforated Dashed Tear Line */}
      <div
        className="flex-1 h-full flex items-center px-1"
        style={{ backgroundColor: fillColor }}
      >
        <div
          className="w-full border-t border-dashed"
          style={{ borderColor: dashedLineColor }}
        />
      </div>

      {/* Right Inward Notch SVG */}
      <svg
        width={width}
        height={height}
        viewBox={`0 0 ${width} ${height}`}
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="shrink-0 -mr-[1px] block"
      >
        {/* Paper Fill inside card (everything left of the concave arc) */}
        <path
          d={`M ${width} 0 L 0 0 L 0 ${height} L ${width} ${height} A ${notchRadius} ${notchRadius} 0 0 1 ${width} 0 Z`}
          fill={fillColor}
        />
        {/* Inward border curve connecting top and bottom card edges */}
        <path
          d={`M ${width - 0.5} 0 A ${notchRadius - 0.5} ${notchRadius - 0.5} 0 0 0 ${width - 0.5} ${height}`}
          stroke={borderColor}
          strokeWidth="1"
          fill="none"
        />
      </svg>
    </div>
  );
}
