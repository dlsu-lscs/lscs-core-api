"use client"

import { useState, useEffect } from "react"
import { api } from "@/lib/api"

function extractObjectKey(imageUrl: string): string {
  if (!imageUrl) return ""
  const idx = imageUrl.indexOf("profile-images/")
  return idx >= 0 ? imageUrl.slice(idx) : ""
}

const urlCache = new Map<string, { url: string; expiresAt: number }>()

interface ProfileImageProps {
  imageUrl?: string | null
  fullName: string
  size?: "sm" | "lg"
  className?: string
}

export function ProfileImage({ imageUrl, fullName, size = "sm", className }: ProfileImageProps) {
  const [src, setSrc] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(false)

  const sizeClass = size === "lg"
    ? "h-24 w-24 text-2xl"
    : "h-8 w-8 text-xs"

  const initials = fullName
    ?.split(" ")
    .filter(Boolean)
    .map((w) => w[0])
    .join("")
    .slice(0, 2)
    .toUpperCase() || "?"

  useEffect(() => {
    if (!imageUrl) return

    // already presigned — use directly
    if (imageUrl.includes("X-Amz-")) {
      setSrc(imageUrl)
      return
    }

    const key = extractObjectKey(imageUrl)
    if (!key) return

    const cached = urlCache.get(key)
    if (cached && cached.expiresAt > Date.now()) {
      setSrc(cached.url)
      return
    }

    setLoading(true)
    setError(false)

    api.getImageUrl(key)
      .then(({ url, expires_in_seconds }) => {
        const ttl = (expires_in_seconds - 120) * 1000 // expire 2 min before server
        urlCache.set(key, { url, expiresAt: Date.now() + ttl })
        setSrc(url)
      })
      .catch(() => setError(true))
      .finally(() => setLoading(false))
  }, [imageUrl])

  if (loading) {
    return (
      <div
        className={`${sizeClass} rounded-full bg-primary/10 animate-pulse shrink-0 ${className ?? ""}`}
      />
    )
  }

  if (error || !src) {
    return (
      <div
        className={`${sizeClass} rounded-full bg-primary/10 flex items-center justify-center shrink-0 font-bold text-primary ${className ?? ""}`}
      >
        <span>{initials}</span>
      </div>
    )
  }

  return (
    <img
      src={src}
      alt={fullName}
      className={`${sizeClass} rounded-full object-cover shrink-0 ${className ?? ""}`}
    />
  )
}
