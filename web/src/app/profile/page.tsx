"use client"

import { useState } from "react"
import { useAuth } from "@/hooks/use-auth"
import { api } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useQueryClient } from "@tanstack/react-query"
import { Loader2, Upload, Trash2, AlertCircle } from "lucide-react"
import { AuthenticatedLayout } from "@/components/authenticated-layout"
import { ProfileImage } from "@/components/profile-image"
import { resizeImage, validateImageSize } from "@/lib/image-utils"

export default function ProfilePage() {
  const { user, updateProfile } = useAuth()
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(false)
  const [formData, setFormData] = useState({
    nickname: user?.nickname || "",
    telegram: user?.telegram || "",
    discord: user?.discord || "",
    interests: user?.interests || "",
    contact_number: user?.contact_number || "",
    fb_link: user?.fb_link || "",
  })
  const [uploading, setUploading] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await updateProfile(formData)
      setIsEditing(false)
    } catch (error) {
      console.error("Failed to update profile:", error)
    }
  }

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    const sizeError = validateImageSize(file)
    if (sizeError) {
      setUploadError(sizeError)
      return
    }

    try {
      setUploading(true)
      setUploadError(null)

      // Resize to max 512px, JPEG quality 0.8
      const resized = await resizeImage(file)

      // Get pre-signed URL (always JPEG after resize)
      const { upload_url, object_key } = await api.generateUploadUrl("image/jpeg")

      // Upload to S3
      await fetch(upload_url, {
        method: "PUT",
        body: resized,
        headers: {
          "Content-Type": "image/jpeg",
        },
      })

      // Complete upload
      await api.completeUpload(object_key)

      // Refresh profile
      queryClient.invalidateQueries({ queryKey: ["member", "me"] })
    } catch (error) {
      console.error("Failed to upload image:", error)
      setUploadError("Upload failed. Please try again.")
    } finally {
      setUploading(false)
    }
   }

  const handleDeleteImage = async () => {
    try {
      await api.deleteImage()
      queryClient.invalidateQueries({ queryKey: ["member", "me"] })
    } catch (error) {
      console.error("Failed to delete image:", error)
    }
  }

  if (!user) {
    return null
  }

  return (
    <AuthenticatedLayout>
      <div className="max-w-2xl mx-auto space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Profile</h1>
        <p className="text-muted-foreground mt-2">
          Manage your profile information
        </p>
      </div>

      {/* Profile Image */}
      <Card>
        <CardHeader>
          <CardTitle>Profile Photo</CardTitle>
          <CardDescription>
            Upload a profile photo to help others recognize you
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-4">
            <ProfileImage imageUrl={user.image_url} fullName={user.full_name} size="lg" />
            <div className="space-y-2">
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => document.getElementById("image-upload")?.click()}
                  disabled={uploading}
                >
                  {uploading ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Upload className="h-4 w-4" />
                  )}
                  <span className="ml-2">Upload Photo</span>
                </Button>
                <Input
                  id="image-upload"
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  className="hidden"
                  onChange={handleImageUpload}
                />
                {user.image_url && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleDeleteImage}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </div>
              <p className="text-xs text-muted-foreground">
                JPG, PNG or WebP. Max 5MB.
              </p>
              {uploadError && (
                <p className="text-xs text-destructive flex items-center gap-1 mt-1">
                  <AlertCircle className="h-3 w-3" />
                  {uploadError}
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Basic Info */}
      <Card>
        <CardHeader>
          <CardTitle>Basic Information</CardTitle>
          <CardDescription>
            Your basic details from the organization database
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="text-sm font-medium">Full Name</label>
              <p className="text-muted-foreground">{user.full_name}</p>
            </div>
            <div>
              <label className="text-sm font-medium">Email</label>
              <p className="text-muted-foreground">{user.email}</p>
            </div>
            <div>
              <label className="text-sm font-medium">Position</label>
              <p className="text-muted-foreground">
                {user.position_name || user.position_id}
              </p>
            </div>
            <div>
              <label className="text-sm font-medium">Committee</label>
              <p className="text-muted-foreground">
                {user.committee_name || user.committee_id}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Editable Info */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Contact & Details</CardTitle>
            <CardDescription>
              Update your contact information and preferences
            </CardDescription>
          </div>
          <Button
            variant={isEditing ? "outline" : "default"}
            size="sm"
            onClick={() => setIsEditing(!isEditing)}
          >
            {isEditing ? "Cancel" : "Edit"}
          </Button>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="text-sm font-medium">Nickname</label>
                <Input
                  value={formData.nickname}
                  onChange={(e) =>
                    setFormData({ ...formData, nickname: e.target.value })
                  }
                  disabled={!isEditing}
                  placeholder="How should we call you?"
                />
              </div>
              <div>
                <label className="text-sm font-medium">Telegram</label>
                <Input
                  value={formData.telegram}
                  onChange={(e) =>
                    setFormData({ ...formData, telegram: e.target.value })
                  }
                  disabled={!isEditing}
                  placeholder="@username"
                />
              </div>
              <div>
                <label className="text-sm font-medium">Discord</label>
                <Input
                  value={formData.discord}
                  onChange={(e) =>
                    setFormData({ ...formData, discord: e.target.value })
                  }
                  disabled={!isEditing}
                  placeholder="username#0000"
                />
              </div>
              <div>
                <label className="text-sm font-medium">Contact Number</label>
                <Input
                  value={formData.contact_number}
                  onChange={(e) =>
                    setFormData({ ...formData, contact_number: e.target.value })
                  }
                  disabled={!isEditing}
                  placeholder="+63..."
                />
              </div>
              <div>
                <label className="text-sm font-medium">Facebook Link</label>
                <Input
                  value={formData.fb_link}
                  onChange={(e) =>
                    setFormData({ ...formData, fb_link: e.target.value })
                  }
                  disabled={!isEditing}
                  placeholder="https://facebook.com/..."
                />
              </div>
            </div>
            <div>
              <label className="text-sm font-medium">Interests</label>
              <Input
                value={formData.interests}
                onChange={(e) =>
                  setFormData({ ...formData, interests: e.target.value })
                }
                disabled={!isEditing}
                placeholder="What are your interests?"
              />
            </div>
            {isEditing && (
              <Button type="submit" className="w-full">
                Save Changes
              </Button>
            )}
          </form>
        </CardContent>
        </Card>
      </div>
    </AuthenticatedLayout>
  )
}
